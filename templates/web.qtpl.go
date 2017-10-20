// This file is automatically generated by qtc from "web.qtpl".
// See https://github.com/valyala/quicktemplate for details.

//line templates/web.qtpl:1
package templates

//line templates/web.qtpl:1
import "github.com/mkreder/dockerpanel/model"

//line templates/web.qtpl:2
import "strconv"

//line templates/web.qtpl:3
import (
	qtio422016 "io"

	qt422016 "github.com/valyala/quicktemplate"
)

//line templates/web.qtpl:3
var (
	_ = qtio422016.Copy
	_ = qt422016.AcquireByteBuffer
)

//line templates/web.qtpl:3
func StreamWebTemplate(qw422016 *qt422016.Writer, webs []model.Web, error string) {
	//line templates/web.qtpl:3
	qw422016.N().S(`

<!DOCTYPE html>
<html lang="en">
<script type="text/javascript">
    function validateForm() {
        var x = document.getElementById("dominio").value;
        if (x == "" ) {
            alert("Se debe completar el dominio");
            return false;
        }

        if (/^[a-zA-Z0-9][a-zA-Z0-9-]{1,61}[a-zA-Z0-9](?:\.[a-zA-Z]{2,})+$/.test(x)) {
        } else {
            alert("Nombre de dominio invalido");
            return false;
        }
    }

    function showDiv() {
        var x = document.getElementById('form');
        x.style.display = 'block';
    }

    function activateSelects(){
        if ( document.getElementById('php').checked ) {
            document.getElementById('phpVersion').removeAttribute("disabled");
            document.getElementById('phpMethod').removeAttribute("disabled");
        } else {
            document.getElementById('phpVersion').setAttribute("disabled","disabled")
            document.getElementById('phpMethod').setAttribute("disabled","disabled");
        }

    }

    function activateFileUpload(){
        if ( document.getElementById('ssl').checked ) {
            document.getElementById('files').removeAttribute("disabled");
        } else {
            document.getElementById('files').setAttribute("disabled","disabled")
        }

    }

    function hideDiv() {
        var x = document.getElementById('form');
        x.style.display = 'none';
        document.getElementById("dominio").value = "";
        document.getElementById("id").value = "";
        document.getElementById("ssl").checked = false;
        document.getElementById("cgi").checked = false;
        document.getElementById("php").checked = false;
        document.getElementById("phpVersion").options[0].selected = true;
        document.getElementById("phpMethod").options[0].selected = true;
        document.getElementById("webserver").options[0].selected = true;
        document.getElementById('phpVersion').setAttribute("disabled","disabled")
        document.getElementById('phpMethod').setAttribute("disabled","disabled");
        document.getElementById('files').setAttribute("disabled","disabled")
    }

    function modifyWeb(id,dominio,ssl,cgi,php,phpVersion,phpMethod,webserver) {
        document.getElementById("id").value = id;
        document.getElementById("dominio").value = dominio;

        if ( ssl == true ){
            document.getElementById("ssl").checked = true;
            document.getElementById('files').removeAttribute("disabled");
        } else {
            document.getElementById("ssl").checked = false;
        }
        if ( cgi == true ) {
            document.getElementById("cgi").checked = true;
        } else {
            document.getElementById("cgi").checked = false;
        }
        if ( php == true ) {
            document.getElementById("php").checked = true;
            document.getElementById('phpVersion').removeAttribute("disabled");
            document.getElementById('phpMethod').removeAttribute("disabled");
        } else {
            document.getElementById("php").checked = false;
        }
        var x = document.getElementById('form');

        var selphpVersion = document.getElementById("phpVersion");


        for (var i = 0; i < selphpVersion.options.length; i++) {
            if (selphpVersion.options[i].text== phpVersion) {
                selphpVersion.options[i].selected = true;
            }
        }

        var selphpMethod = document.getElementById("phpMethod");

        for (var i = 0; i < selphpMethod.options.length; i++) {
            if (selphpMethod.options[i].text== phpMethod) {
                selphpMethod.options[i].selected = true;
            }
        }

        var selwebserver = document.getElementById("webserver");

        for (var i = 0; i < selwebserver.options.length; i++) {
            if (selwebserver.options[i].text== webserver) {
                selwebserver.options[i].selected = true;
            }
        }


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
	//line templates/web.qtpl:153
	if len(error) > 0 {
		//line templates/web.qtpl:153
		qw422016.N().S(`
<script type="text/javascript">
    window.alert("`)
		//line templates/web.qtpl:155
		qw422016.E().S(error)
		//line templates/web.qtpl:155
		qw422016.N().S(`")
</script>
`)
		//line templates/web.qtpl:157
	}
	//line templates/web.qtpl:157
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
                        <h4 class="panel-title pull-left" style="padding-top: 7.5px;">Sitios Web</h4>
                        <button type="button" class="btn pull-right btn-primary btn-sm" onclick="showDiv()">Agregar</button>
                    </div>
                    <!-- /.panel-heading -->
                    <div class="panel-body">
                        <table width="100%" class="table table-striped table-bordered table-hover" id="dataTables-example">
                            <thead>
                            <tr>
                                <th>Dominio</th>
                                <th>Estado</th>
                                <th>Servidor Web</th>
                                <th>SSL</th>
                                <th>PHP</th>
                                <th>CGI</th>
                                <th>Acciones</th>
                            </tr>
                            </thead>
                            <tbody>
                            `)
	//line templates/web.qtpl:245
	for _, web := range webs {
		//line templates/web.qtpl:245
		qw422016.N().S(`
                            <tr class="odd gradeX">
                                <td> `)
		//line templates/web.qtpl:247
		qw422016.E().S(web.Dominio)
		//line templates/web.qtpl:247
		qw422016.N().S(` </td>


                                `)
		//line templates/web.qtpl:250
		switch web.Estado {
		//line templates/web.qtpl:251
		case 1:
			//line templates/web.qtpl:251
			qw422016.N().S(`
                                <td>a configurar</td>
                                `)
		//line templates/web.qtpl:253
		case 2:
			//line templates/web.qtpl:253
			qw422016.N().S(`
                                <td>configurado</td>
                                `)
		//line templates/web.qtpl:255
		case 3:
			//line templates/web.qtpl:255
			qw422016.N().S(`
                                <td>corriendo</td>
                                `)
			//line templates/web.qtpl:257
		}
		//line templates/web.qtpl:257
		qw422016.N().S(`

                                <td> `)
		//line templates/web.qtpl:259
		qw422016.E().S(web.Webserver)
		//line templates/web.qtpl:259
		qw422016.N().S(` </td>


                                `)
		//line templates/web.qtpl:262
		if web.SSL {
			//line templates/web.qtpl:262
			qw422016.N().S(`
                                <td class="center">si</td>
                                `)
			//line templates/web.qtpl:264
		} else {
			//line templates/web.qtpl:264
			qw422016.N().S(`
                                <td class="center">no</td>
                                `)
			//line templates/web.qtpl:266
		}
		//line templates/web.qtpl:266
		qw422016.N().S(`

                                `)
		//line templates/web.qtpl:268
		if web.PHP {
			//line templates/web.qtpl:268
			qw422016.N().S(`
                                <td class="center">si</td>
                                `)
			//line templates/web.qtpl:270
		} else {
			//line templates/web.qtpl:270
			qw422016.N().S(`
                                <td class="center">no</td>
                                `)
			//line templates/web.qtpl:272
		}
		//line templates/web.qtpl:272
		qw422016.N().S(`

                                `)
		//line templates/web.qtpl:274
		if web.CGI {
			//line templates/web.qtpl:274
			qw422016.N().S(`
                                <td class="center">si</td>
                                `)
			//line templates/web.qtpl:276
		} else {
			//line templates/web.qtpl:276
			qw422016.N().S(`
                                <td class="center">no</td>
                                `)
			//line templates/web.qtpl:278
		}
		//line templates/web.qtpl:278
		qw422016.N().S(`



                                <td class="center">
                                    <button type="button" class="btn btn-xs btn-primary" data-toggle="tooltip" data-placement="top" title="Editar sitio Web" onclick='modifyWeb(`)
		//line templates/web.qtpl:283
		qw422016.N().D(int(web.ID))
		//line templates/web.qtpl:283
		qw422016.N().S(`, "`)
		//line templates/web.qtpl:283
		qw422016.E().S(web.Dominio)
		//line templates/web.qtpl:283
		qw422016.N().S(`", `)
		//line templates/web.qtpl:283
		qw422016.E().S(strconv.FormatBool(web.SSL))
		//line templates/web.qtpl:283
		qw422016.N().S(`, `)
		//line templates/web.qtpl:283
		qw422016.E().S(strconv.FormatBool(web.CGI))
		//line templates/web.qtpl:283
		qw422016.N().S(`, `)
		//line templates/web.qtpl:283
		qw422016.E().S(strconv.FormatBool(web.PHP))
		//line templates/web.qtpl:283
		qw422016.N().S(`, "`)
		//line templates/web.qtpl:283
		qw422016.E().S(web.PHPversion)
		//line templates/web.qtpl:283
		qw422016.N().S(`", "`)
		//line templates/web.qtpl:283
		qw422016.E().S(web.PHPmethod)
		//line templates/web.qtpl:283
		qw422016.N().S(`" , "`)
		//line templates/web.qtpl:283
		qw422016.E().S(web.Webserver)
		//line templates/web.qtpl:283
		qw422016.N().S(`" )' ><i class="fa fa-list"></i></button>
                                    <button class="btn btn-xs btn-danger" data-toggle="tooltip" data-placement="top" title="Eliminar sitio Web" onclick="location.href='removeWeb?id=`)
		//line templates/web.qtpl:284
		qw422016.N().D(int(web.ID))
		//line templates/web.qtpl:284
		qw422016.N().S(`';"><i class="fa fa-trash-o"></i></button>
                                </td>
                            </tr>
                            `)
		//line templates/web.qtpl:287
	}
	//line templates/web.qtpl:287
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
                    Configuración del Sitio
                </div>
                <div class="panel-body">
                    <form id="addweb" action="/web" enctype="multipart/form-data" onsubmit="return validateForm()" role=form method="post">
                        <label>Dominio</label>
                        <input id="id" name="id" hidden="true" >
                        <input class="form-control" name="dominio" id="dominio">
                        <br>
                        <label>Configuración</label>
                        <br>
                        Servidor Web
                        <select  class="form-control"  id="webserver" name="webserver">
                            <option value="apache">Apache</option>
                            <option value="nginx">Nginx</option>
                        </select>
                        <div class="checkbox">
                            <label>
                                <input id="cgi" name="cgi" type="checkbox" value="true"> CGI
                            </label>
                        </div>
                        <div class="checkbox">
                            <label>
                                <input id="php" name="php" type="checkbox" value="true" onclick="activateSelects()"> PHP
                            </label>
                        </div>
                        <div class="checkbox">
                        <label>
                            <input id="ssl" name="ssl" type="checkbox" value="true" onclick="activateFileUpload()"> SSL
                        </label>
                    </div>
                        <label for="files" class="btn btn-default btn-xs">Subir certificado y llave SSL</label>
                        <input id="files" style="visibility:hidden;" type="file" name="pem" disabled="disabled">

                        Versión de PHP
                        <select  class="form-control"  id="phpVersion" name="phpVersion" disabled>
                            <option value="5.4">5.4</option>
                            <option value="5.6">5.6</option>
                            <option value="7.0">7.0</option>
                            <option value="7.1">7.1</option>
                        </select>

                        <br>
                        Metodo PHP
                        <select   class="form-control" id="phpMethod" name="phpMethod" disabled=>
                            <option value="fpm">PHP-FPM</option>
                            <option value="fcgi">FastCGI</option>
                        </select>

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
//line templates/web.qtpl:384
}

//line templates/web.qtpl:384
func WriteWebTemplate(qq422016 qtio422016.Writer, webs []model.Web, error string) {
	//line templates/web.qtpl:384
	qw422016 := qt422016.AcquireWriter(qq422016)
	//line templates/web.qtpl:384
	StreamWebTemplate(qw422016, webs, error)
	//line templates/web.qtpl:384
	qt422016.ReleaseWriter(qw422016)
//line templates/web.qtpl:384
}

//line templates/web.qtpl:384
func WebTemplate(webs []model.Web, error string) string {
	//line templates/web.qtpl:384
	qb422016 := qt422016.AcquireByteBuffer()
	//line templates/web.qtpl:384
	WriteWebTemplate(qb422016, webs, error)
	//line templates/web.qtpl:384
	qs422016 := string(qb422016.B)
	//line templates/web.qtpl:384
	qt422016.ReleaseByteBuffer(qb422016)
	//line templates/web.qtpl:384
	return qs422016
//line templates/web.qtpl:384
}
