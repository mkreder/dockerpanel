package worker

import (
	"context"
	"github.com/mkreder/dockerpanel/model"
	"github.com/docker/docker/client"
	"github.com/docker/docker/api/types"
	"github.com/jhoonb/archivex"
	"log"
	"os"
	"os/exec"
	"time"

	"strings"
	"strconv"
	"io/ioutil"
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

				if _ , err := os.Stat("configs/dns/" + zona.Dominio + ".conf"); os.IsNotExist(err) {
					f, err := os.Create("configs/dns/" + zona.Dominio + ".conf")
					check(err)
					defer f.Close()

					_, _ = f.WriteString("zone " + zona.Dominio + " IN { \n" +
						"   type master \n" +
						"    file \"/etc/bind/zones/" + zona.Dominio + ".conf\"; \n" +
						"};\n")
					f.Sync()

					dockerf, err := os.OpenFile("configs/dns/Dockerfile", os.O_APPEND|os.O_WRONLY, 0600)
					if err != nil {
						panic(err)
					}

					defer dockerf.Close()

					dockerf.WriteString("COPY " + zona.Dominio + ".conf /etc/bind/conf.d/ \n" +
						"COPY zone-" + zona.Dominio + ".conf /etc/bind/zones/ \n")

				}

				nsname := ""
				for _ , registro := range zona.Registros {
					if ( registro.Tipo == "NS" ){
						nsname = registro.Valor
						break
					}
				}

				serial := strconv.Itoa(int(time.Now().Unix()))
				nuevaZona := "@	IN	SOA " + nsname + "\t"+ strings.Replace(zona.Email,"@",".",-1)+ " ( \n"+
					"\t\t\t\t\t\t" + serial + "\n" +
					"\t\t\t\t\t\t 8H \n" +
					"\t\t\t\t\t\t 2H \n" +
					"\t\t\t\t\t\t 4W \n" +
					"\t\t\t\t\t\t 1D ) \n"

				for _ , registro := range zona.Registros {
					if ( registro.Tipo == "NS" ){
						nuevaZona = nuevaZona + "\t\t\t\t\t\t\t NS \t" + registro.Valor + "\n"
					}
				}

				for _ , registro := range zona.Registros {
					if ( registro.Tipo == "MX" ){
						nuevaZona = nuevaZona + "\t\t\t\t\t\t\t MX \t" + registro.Prioridad + " " + registro.Valor + "\n"
					}
				}

				for _ , registro := range zona.Registros {
					if ( registro.Tipo != "MX" ) && ( registro.Tipo != "NS" ) {
						nuevaZona = nuevaZona + registro.Nombre + "\t\t" + registro.Tipo + "\t\t" + registro.Valor + "\n"
					}
				}

				f2, err := os.Create("configs/dns/zone-" + zona.Dominio + ".conf" )
				check(err)
				defer f2.Close()

				_, _ = f2.WriteString(nuevaZona)
				f2.Sync()

				os.MkdirAll("configs/container",0755)
				tar := new(archivex.TarFile)
				tar.Create("configs/container/dnsconf.tar")
				tar.AddAll("configs/dns", false)
				tar.Close()

				dockerBuildContext, err := os.Open("config/dns/dnsconf.tar")
				defer dockerBuildContext.Close()


				buildOptions := types.ImageBuildOptions{
					CPUSetCPUs:   "2",
					CPUSetMems:   "12",
					CPUShares:    20,
					CPUQuota:     10,
					CPUPeriod:    30,
					Memory:       256,
					MemorySwap:   512,
					ShmSize:      10,
					CgroupParent: "cgroup_parent",
					Dockerfile:   "Dockerfile", // optional, is the default
					Tags:   []string{"dnsimage"},
				}

				cli, err := client.NewEnvClient()

				buildResponse, err := cli.ImageBuild(context.Background(), dockerBuildContext, buildOptions)
				if err != nil {
					log.Printf("%s", err.Error())
				}
				if err != nil {
					log.Fatal(err)
				}
				response, err := ioutil.ReadAll(buildResponse.Body)
				log.Printf(string(response))

				defer buildResponse.Body.Close()


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