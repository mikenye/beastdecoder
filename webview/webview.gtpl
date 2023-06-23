<!DOCTYPE html>
<html>
  <head>
    <style>
      table, th, td {
        border: 1px solid black;
        border-collapse: collapse;
      }
      body {
        font-family: "Lucida Console", "Courier New", monospace;
      }
    </style>
    <meta http-equiv="Refresh" content="1"> 
  </head>
  <body>
    <table style="width:100%">
      <tr>
        <th>ICAO</th>
        <th>Sqwk</th>
        <th>Call</th>
        <th>Alt</th>
        <th>Lat</th>
        <th>Lon</th>
        <th>Method</th>
        <th>Spd</th>
        <th>Hdg</th>
        <th>Msgs</th>
      </tr>
    {{range $index, $element := .}}
      <tr>
        <td>{{printf "%06x" $index}}</td>
        <td>
          {{if .SquawkCodeKnown}}
            {{.SquawkCode}}
          {{else}}
            &nbsp;
          {{end}}
        </td>
        <td>{{if .CallsignKnown}}{{.Callsign}}{{else}}&nbsp;{{end}}</td>
        <td>
          {{if .AirborneStatusKnown}}
            {{if .Airborne}}
              {{ if .AltitudeKnown}}
                {{.Altitude}}
              {{else}}
                &nbsp;
              {{end}}
            {{else}}
              ground
            {{end}}
          {{end}}
        </td>
        <td>
          {{if .LatLonKnown}}
            {{printf "%.5f" .Lat}}
          {{end}}
        </td>
        <td>
          {{if .LatLonKnown}}
            {{printf "%.5f" .Lon}}
          {{end}}
        </td>
        <td>
          {{if .LatLonKnown}}
            {{.LatLonMethod}}
          {{end}}
        </td>
        <td>
          {{if .GroundSpeedKnown}}
            {{.GroundSpeed}}
          {{end}}
        </td>
        <td>
          {{if .GroundTrackKnown}}
            {{.GroundTrack}}
          {{end}}
        </td>
        <td>
          {{.MsgCount}}
        </td>
      </tr>
    {{end}}
    </table>
  </body>
</html>
