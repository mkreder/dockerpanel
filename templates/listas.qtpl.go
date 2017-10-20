// This file is automatically generated by qtc from "listas.qtpl".
// See https://github.com/valyala/quicktemplate for details.

//line templates/listas.qtpl:1
package templates

//line templates/listas.qtpl:1
import "github.com/mkreder/dockerpanel/model"

//line templates/listas.qtpl:2
import (
	qtio422016 "io"

	qt422016 "github.com/valyala/quicktemplate"
)

//line templates/listas.qtpl:2
var (
	_ = qtio422016.Copy
	_ = qt422016.AcquireByteBuffer
)

//line templates/listas.qtpl:2
func StreamListaTemplate(qw422016 *qt422016.Writer, listas []model.Lista, dominioid string, error string) {
	//line templates/listas.qtpl:2
	qw422016.N().S(`

<!DOCTYPE html>
<html lang="en">
<script type="text/javascript">
    function validateForm() {
        var y = document.getElementById("nombre").value;
        if (y == "" ) {
            alert("Se debe completar el nombre de la lista");
            return false;
        }

        var z = document.getElementById("password").value;
        if ( document.getElementById("password").getAttribute("disabled") == "disabled" ){
            } else {
            if (z == "" ) {
                alert("Se debe completar la contraseña");
                return false;
            }
        }

    }

    function generarPassword() {
        if ( document.getElementById("password").getAttribute("disabled") == "disabled" ) {
            document.getElementById("btngenerar").innerHTML="Generar";
            document.getElementById("password").removeAttribute("disabled")
            document.getElementById("checkmostrar").removeAttribute("disabled");
        } else {
            document.getElementById("password").value = Math.random().toString(36).slice(-8);
        }
    }

    function mostrarPassword() {
        if ( document.getElementById('checkmostrar').checked ) {
            document.getElementById('password').removeAttribute("type")
        } else {
            document.getElementById('password').setAttribute("type","password")
        }

    }


    function showDiv() {
        var x = document.getElementById('form');
        x.style.display = 'block';
    }


    function hideDiv() {
        var x = document.getElementById('form');
        x.style.display = 'none';
        document.getElementById("nombre").value = "";
        document.getElementById("password").value = "";
        if ( document.getElementById("password").getAttribute("disabled") == "disabled" ) {
            document.getElementById("btngenerar").innerHTML="Generar";
            document.getElementById("password").removeAttribute("disabled");
            document.getElementById("checkmostrar").removeAttribute("disabled");
        }
    }

    function modifyLista(id,nombre) {
        var x = document.getElementById('form');
        document.getElementById("id").value = id;
        document.getElementById("nombre").value = nombre;
        document.getElementById("password").setAttribute("disabled","disabled")
        document.getElementById("btngenerar").innerHTML="Cambiar";
        document.getElementById("checkmostrar").setAttribute("disabled","disabled")
        x.style.display = 'block';
    }

</script>
<head>

    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <meta name="description" content="">
    <meta name="author" content="">

    <title>Panel Docker</title>

    <!-- Bootstrap Core CSS -->
    <link href="../vendor/bootstrap/css/bootstrap.min.css" rel="stylesheet">

    <!-- MetisMenu CSS -->
    <link href="../vendor/metisMenu/metisMenu.min.css" rel="stylesheet">

    <!-- Custom CSS -->
    <link href="../dist/css/sb-admin-2.css" rel="stylesheet">

    <!-- Morris Charts CSS -->
    <link href="../vendor/morrisjs/morris.css" rel="stylesheet">

    <!-- Custom Fonts -->
    <link href="../vendor/font-awesome/css/font-awesome.min.css" rel="stylesheet" type="text/css">

    <!-- HTML5 Shim and Respond.js IE8 support of HTML5 elements and media queries -->
    <!-- WARNING: Respond.js doesn't work if you view the page via file:// -->
    <!--[if lt IE 9]>
    <script src="https://oss.maxcdn.com/libs/html5shiv/3.7.0/html5shiv.js"></script>
    <script src="https://oss.maxcdn.com/libs/respond.js/1.4.2/respond.min.js"></script>
    <![endif]-->

</head>

<body>

`)
	//line templates/listas.qtpl:110
	if len(error) > 0 {
		//line templates/listas.qtpl:110
		qw422016.N().S(`
<script type="text/javascript">
    window.alert("`)
		//line templates/listas.qtpl:112
		qw422016.E().S(error)
		//line templates/listas.qtpl:112
		qw422016.N().S(`")
</script>
`)
		//line templates/listas.qtpl:114
	}
	//line templates/listas.qtpl:114
	qw422016.N().S(`

<div id="wrapper">

    <!-- Navigation -->
    <nav class="navbar navbar-default navbar-static-top" role="navigation" style="margin-bottom: 0">
        <div class="navbar-header">
            <button type="button" class="navbar-toggle" data-toggle="collapse" data-target=".navbar-collapse">
                <span class="sr-only">Toggle navigation</span>
                <span class="icon-bar"></span>
                <span class="icon-bar"></span>
                <span class="icon-bar"></span>
            </button>
            <a class="navbar-brand" href="index.html">Docker Panel</a>
        </div>
        <!-- /.navbar-header -->

        <ul class="nav navbar-top-links navbar-right">
            <!-- /.dropdown -->
            <li class="dropdown">
                <a class="dropdown-toggle" data-toggle="dropdown" href="#">
                    <i class="fa fa-user fa-fw"></i> <i class="fa fa-caret-down"></i>
                </a>
                <ul class="dropdown-menu dropdown-user">
                    <li><a href="/profile"><i class="fa fa-user fa-fw"></i> Configuración</a>
                    </li>
                    <li class="divider"></li>
                    <li><a href="/logout"><i class="fa fa-sign-out fa-fw"></i> Logout</a>
                    </li>
                </ul>
                <!-- /.dropdown-user -->
            </li>
            <!-- /.dropdown -->
        </ul>
        <!-- /.navbar-top-links -->

        <div class="navbar-default sidebar" role="navigation">
            <div class="sidebar-nav navbar-collapse">
                <ul class="nav" id="side-menu">
                    <li>
                        <a href="/"><i class="fa fa-dashboard fa-fw"></i> Dashboard</a>
                    </li>
                    <li>
                        <a href="/web"><i class="fa fa-server fa-fw"></i>Sitios Web</a>
                    </li>
                    <li>
                        <a href="/dns"><i class="fa fa-cloud fa-fw"></i>DNS</a>
                    </li>
                    <li>
                        <a href="/bd"><i class="fa fa-database fa-fw"></i>Base de Datos</a>
                    </li>
                    <li>
                        <a href="/mail"><i class="fa fa-at fa-fw"></i>E-Mail</a>
                    </li>
                    <li>
                        <a href="/ftp"><i class="fa fa-file-archive-o fa-fw"></i>FTP</a>
                    </li>
                </ul>
            </div>
            <!-- /.sidebar-collapse -->
        </div>
        <!-- /.navbar-static-side -->
    </nav>

    <div id="page-wrapper">
        <br>
        <div class="row">
            <div class="col-lg-12">
                <div class="panel panel-default">
                    <div class="panel-heading clearfix">
                        <h4 class="panel-title pull-left" style="padding-top: 7.5px;">Listas de correo</h4>
                        <button type="button" class="btn pull-right btn-primary btn-sm" onclick="showDiv()">Agregar</button>
                    </div>
                    <!-- /.panel-heading -->
                    <div class="panel-body">
                        <table width="100%" class="table table-striped table-bordered table-hover" id="dataTables-example">
                            <thead>
                            <tr>
                                <th>Nombre</th>
                                <th>Estado</th>
                                <th>Acciones</th>
                            </tr>
                            </thead>
                            <tbody>
                            `)
	//line templates/listas.qtpl:198
	for _, lista := range listas {
		//line templates/listas.qtpl:198
		qw422016.N().S(`
                            <tr class="odd gradeX">
                                <td> `)
		//line templates/listas.qtpl:200
		qw422016.E().S(lista.Nombre)
		//line templates/listas.qtpl:200
		qw422016.N().S(` </td>

                                `)
		//line templates/listas.qtpl:202
		switch lista.Estado {
		//line templates/listas.qtpl:203
		case 1:
			//line templates/listas.qtpl:203
			qw422016.N().S(`
                                <td>a configurar</td>
                                `)
		//line templates/listas.qtpl:205
		case 2:
			//line templates/listas.qtpl:205
			qw422016.N().S(`
                                <td>configurado</td>
                                `)
			//line templates/listas.qtpl:207
		}
		//line templates/listas.qtpl:207
		qw422016.N().S(`


                                <td class="center">
                                    <button type="button" class="btn btn-xs btn-primary" data-toggle="tooltip" data-placement="top" title="Modificar lista de correo" onclick='modifyLista(`)
		//line templates/listas.qtpl:211
		qw422016.N().D(int(lista.ID))
		//line templates/listas.qtpl:211
		qw422016.N().S(`, "`)
		//line templates/listas.qtpl:211
		qw422016.E().S(lista.Nombre)
		//line templates/listas.qtpl:211
		qw422016.N().S(`")' ><i class="fa fa-list"></i></button>
                                    <button class="btn btn-xs btn-danger" data-toggle="tooltip" data-placement="top" title="Eliminar lista de correo" onclick="location.href='removeLista?id=`)
		//line templates/listas.qtpl:212
		qw422016.N().D(int(lista.ID))
		//line templates/listas.qtpl:212
		qw422016.N().S(`&dominioid=`)
		//line templates/listas.qtpl:212
		qw422016.E().S(dominioid)
		//line templates/listas.qtpl:212
		qw422016.N().S(`';"><i class="fa fa-trash-o"></i></button>
                                </td>
                            </tr>
                            `)
		//line templates/listas.qtpl:215
	}
	//line templates/listas.qtpl:215
	qw422016.N().S(`
                            </tbody>
                        </table>
                    </div>
                    <!-- /.panel-body -->
                </div>
                <!-- /.panel -->
            </div>
            <!-- /.col-lg-12 -->
        </div>


        <div id="form" class="col-lg-6" hidden="true" >
            <div class="panel panel-default">
                <div class="panel-heading">
                    Configuración de la lista de Correo
                </div>
                <div class="panel-body">
                    <form id="addftp" action="/addLista" onsubmit="return validateForm()" role=form method="post">
                        <input id="id" name="id" hidden="true" >
                        <input id="dominioid" name="dominioid" hidden="true" value="`)
	//line templates/listas.qtpl:235
	qw422016.E().S(dominioid)
	//line templates/listas.qtpl:235
	qw422016.N().S(`">
                        <label>Configuración</label>
                        <br>
                        Nombre de la Lista
                        <input class="form-control" name="nombre" id="nombre">
                        <br>

                        Contraseña
                        <input class="form-control" name="password" id="password" type="password">
                        <button id="btngenerar" type="button" class="btn btn-default btn-sm" onclick="generarPassword()">Generar</button>
                        <input id="checkmostrar" name="checkmostrar" type="checkbox" value="true" onclick="mostrarPassword()"> Mostrar Contraseña
                        <br>
                        <br>
                        <button type="submit" class="btn btn-default">Guardar</button>
                        <button type="button" class="btn btn-default" onclick="hideDiv()">Cancelar</button>
                    </form>
                </div>
            </div>

        </div>
        <!-- /#page-wrapper -->

    </div>
    <!-- /#wrapper -->

    <!-- jQuery -->
    <script src="../vendor/jquery/jquery.min.js"></script>

    <!-- Bootstrap Core JavaScript -->
    <script src="../vendor/bootstrap/js/bootstrap.min.js"></script>

    <!-- Metis Menu Plugin JavaScript -->
    <script src="../vendor/metisMenu/metisMenu.min.js"></script>

    <!-- Morris Charts JavaScript -->
    <!--<script src="../vendor/raphael/raphael.min.js"></script>-->
    <!--<script src="../vendor/morrisjs/morris.min.js"></script>-->
    <!--<script src="../data/morris-data.js"></script>-->

    <!-- Custom Theme JavaScript -->
    <script src="../dist/js/sb-admin-2.js"></script>

</body>

</html>

`)
//line templates/listas.qtpl:281
}

//line templates/listas.qtpl:281
func WriteListaTemplate(qq422016 qtio422016.Writer, listas []model.Lista, dominioid string, error string) {
	//line templates/listas.qtpl:281
	qw422016 := qt422016.AcquireWriter(qq422016)
	//line templates/listas.qtpl:281
	StreamListaTemplate(qw422016, listas, dominioid, error)
	//line templates/listas.qtpl:281
	qt422016.ReleaseWriter(qw422016)
//line templates/listas.qtpl:281
}

//line templates/listas.qtpl:281
func ListaTemplate(listas []model.Lista, dominioid string, error string) string {
	//line templates/listas.qtpl:281
	qb422016 := qt422016.AcquireByteBuffer()
	//line templates/listas.qtpl:281
	WriteListaTemplate(qb422016, listas, dominioid, error)
	//line templates/listas.qtpl:281
	qs422016 := string(qb422016.B)
	//line templates/listas.qtpl:281
	qt422016.ReleaseByteBuffer(qb422016)
	//line templates/listas.qtpl:281
	return qs422016
//line templates/listas.qtpl:281
}