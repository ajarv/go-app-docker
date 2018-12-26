<!DOCTYPE html>
<html lang="en">

<head>
  <title>{{.PageTitle}}</title>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/4.1.3/css/bootstrap.min.css">
  <script src="https://ajax.googleapis.com/ajax/libs/jquery/3.3.1/jquery.min.js"></script>
  <script src="https://cdnjs.cloudflare.com/ajax/libs/popper.js/1.14.3/umd/popper.min.js"></script>
  <script src="https://maxcdn.bootstrapcdn.com/bootstrap/4.1.3/js/bootstrap.min.js"></script>
  <link href="https://fonts.googleapis.com/css?family=Didact+Gothic|Ubuntu" rel="stylesheet">
  <link rel="shortcut icon" type="image/png" href="/static/batman.png" />

</head>
<style>
  body {
    font-family: 'Didact Gothic', sans-serif;
    max-width: 1200px;
    margin: auto;
  }

  .debuginfo {
    font-family: 'Ubuntu', sans-serif;
  }

  .notimportant {
    color: white;
    background-color: aliceblue;
  }

  .sayings {
    font-size: x-small;
  }
</style>

<body>

  <div class="jumbotron text-info">
    <h2 class="notimportant">All important things in life start with a white Screen</h2>
    <h4 class="text-right">
      -- <small>
        {{.AppName}}.{{.ApiVersion}} | Host:{{.Host}}
        {{ if .Environment.CF_INSTANCE_INDEX }}
        [ {{.Environment.CF_INSTANCE_INDEX}} ]
        {{end}}
      </small> </h4>
    {{ if .message }}
    <p class="alert alert-info">{{.message}}</p>
    {{end}}
  </div>


  <div class="container sayings notimportant">
    <div class="row">
      <div class="col-sm-4">
        <h4>Remember there are two type of addictions</h4>
        <p></p>
        <p> The first where euphoria comes first followed by a LOT of questioning.<br />
          .. friday night , party .. you drink yourself to brim and Saturday morning you wake up in hangover and pain
          questioning how did you end up like this.
        </p>
        <p>The second type where questioning comes first and euphoria follows.<br />
          .. 5am .. you are outside running in the cold, questioning how did you end up on the street like this. At 7
          am
          you are back in your home , drinking coffee .. pure euphoria.
        </p>
        <br />
        <h3>Choose your addictions wisely!</h3>
      </div>
      <div class="col-sm-8">
        <div>
          <img alt="" src="/static/batman.png" style="max-width:300px"><br />
        </div>
        <br />

      </div>
    </div>
  </div>

  <div class="container">

    <div class="card font-small">
      <div class="card-header">
        <h2>Request headers</h2>
      </div>
      <div class="card-body debuginfo">
        <table class="table table-sm">
          <thead>
            <tr>
              <th scope="col">Name</th>
              <th scope="col">Value</th>
            </tr>
          </thead>
          <tbody>
            {{ range $key, $value := .Request.Headers }}
            <tr>
              <td>{{$key}}</td>
              <td>
                <span class="d-inline-block text-truncate" style="max-width: 700px;">
                  {{$value}}
                </span>
              </td>
            </tr>
            {{ end }}
          </tbody>
        </table>
      </div>
    </div>
  </div>
  <div class="container">

    {{if .Environment}}
    <div class="card font-small">
      <div class="card-header">
        <h2>Environment Vars</h2>
      </div>
      <div class="card-body debuginfo">
        <table class="table table-sm">
          <thead>
            <tr>
              <th scope="col">Name</th>
              <th scope="col">Value</th>
            </tr>
          </thead>
          <tbody>
            {{ range $key, $value := .Environment }}
            <tr>
              <td>{{$key}}</td>
              <td>
                <span class="d-inline-block text-truncate" style="max-width: 700px;">
                  {{$value}}
                </span>
              </td>
            </tr>
            {{ end }}
          </tbody>
        </table>
      </div>
    </div>
    {{ end }}
  </div>
  <footer class="page-footer font-small blue pt-4">
    <div class="footer-copyright text-center py-3">GitHub:
      <a href="https://github.com/ajarv/go-app-docker.git"> https://github.com/ajarv/go-app-docker.git</a>
    </div>
  </footer>
</body>

</html>