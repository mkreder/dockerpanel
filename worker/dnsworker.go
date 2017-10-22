package worker

import (
	//"context"
	//"github.com/docker/docker/client"
	"github.com/mkreder/dockerpanel/model"
	//"github.com/docker/docker/api/types"
	//"github.com/docker/docker/api/types/container"
	"log"
	"os"
	"os/exec"
	"time"

	"strings"
	"strconv"
)

func check(e error){
	if e != nil {
		panic(e)
	}
}

func RunDNSWorker (){
	// Loop para siempre
	for {
		for _ , zona := range model.Mgr.GetAllZonas() {
			if zona.Estado == 1 {
				log.Printf("Trabajando en la zona " + zona.Dominio )

				if _ , err := os.Stat("configs/dns/Dockerfile"); os.IsNotExist(err){
					// No hay un Dockerfile creado
					srcFolder := "defaults/dns"
					os.MkdirAll("configs",0755)
					destFolder := "configs/"
					cpCmd := exec.Command("cp", "-rf", srcFolder, destFolder)
					err := cpCmd.Run()
					check(err)
				}

				f, err := os.Create("configs/dns/" + zona.Dominio + ".conf" )
				check(err)
				defer f.Close()

				_, _ = f.WriteString("zone " + zona.Dominio + " IN { \n" +
					"   type master \n" +
						"    file \"/etc/bind/zones/" + zona.Dominio + ".conf\"; \n" +
							"};\n")
				f.Sync()

				dockerf, err := os.OpenFile("configs/dns/Dockerfile",os.O_APPEND|os.O_WRONLY,0600)
				if err != nil {
					panic(err)
				}

				defer dockerf.Close()

				dockerf.WriteString("COPY " + zona.Dominio + ".conf /etc/bind/conf.d/ \n" +
					"COPY zone-" + zona.Dominio + ".conf /etc/bind/zones/ \n")


				nsname := ""
				for _ , registro := range zona.Registros {
					if ( registro.Tipo == "NS" ){
						nsname = registro.Valor
						break
					}
				}

				serial := strconv.Itoa(int(time.Now().Unix()))
				nuevaZona := "@	IN	SOA " + nsname + "	"+ strings.Replace(zona.Email,"@",".",-1)+ " ( \n"+
					"			" + serial + "\n" +
					"			8H \n" +
					"			2H \n" +
					"			4W \n" +
					"			1D ) \n"

				for _ , registro := range zona.Registros {
					if ( registro.Tipo == "NS" ){
						nuevaZona = nuevaZona + "			NS	" + registro.Valor + "\n"
					}
				}

				for _ , registro := range zona.Registros {
					if ( registro.Tipo == "MX" ){
						nuevaZona = nuevaZona + "			MX	" + registro.Prioridad + " " + registro.Valor + "\n"
					}
				}

				for _ , registro := range zona.Registros {
					if ( registro.Tipo != "MX" ) && ( registro.Tipo != "NS" ) {
						nuevaZona = nuevaZona + registro.Nombre + "	" + registro.Tipo + "	" + registro.Valor + "\n"
					}
				}
				//TODO desactivar prioridad solo para MX

				f2, err := os.Create("configs/dns/zone-" + zona.Dominio + ".conf" )
				check(err)
				defer f2.Close()

				_, _ = f2.WriteString(nuevaZona)
				f2.Sync()





				zona.Estado = 2
				model.Mgr.UpdateZona(&zona)
			}



			}
		}
		//ctx := context.Background()
		//cli, err := client.NewEnvClient()
		//check(err)
		//
		//
		//_, err = cli.ImagePull(ctx, "docker.io/library/alpine", types.ImagePullOptions{})
		//check(err)
		//
		//cli.ImageCreate(ctx)
		//
		//resp, err := cli.ContainerCreate(ctx, &container.Config{
		//	Image: "alpine",
		//	Cmd:   []string{"echo", "hello world"},
		//}, nil, nil, "")
		//check(err)
		//
		//



}