## About

<table>
<tr>
<td>

This API provide you ready to use v2ray nodes in various format, include your desired cdn/sni host.  
All nodes based on [LatinaSub](https://github.com/LalatinaHub/LatinaSub) result.

</td>
</tr>
</table>

## Endpoints

All endpoints is described at root of the api.  
Check the [deployed API host](https://fool.azurewebsites.net/) as example.

## How To Deploy

- Download and run [LatinaApi docker images](https://github.com/LalatinaHub/LatinaApi/pkgs/container/latinaapi)
  - Default web server port is 8080

## FAQ

Q: Do i need to host/deploy [LatinaSub](https://github.com/LalatinaHub/LatinaSub) too to make an endpoint ?  
A: No, LatinaApi is enough to make an endpoint. By default it will get the database from [LatinaSub](https://github.com/LalatinaHub/LatinaSub), but you can change it if you want.

Q: I have deployed [LatinaSub](https://github.com/LalatinaHub/LatinaSub) and i want to use the result as resources for my LatinaApi, how can i do ?  
A: Change `DbUrl` variable inside `/internal/db/db.go` file

## License

This software released under [MIT License](https://github.com/LalatinaHub/License/blob/main/License)
