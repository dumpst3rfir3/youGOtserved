# You GOt Served

![](OIG3.QmdBp.jpg)

This tool was developed for red team or other security testing purposes. It will simply take a shellcode (e.g., .bin) file, obfuscate the shellcode using [Babble](https://github.com/mjwhitta/babble), and then build a Windows service executable that can run the shellcode. The service executable can also take arguments to install, start, stop, or remove the Windows service.

When the service executable runs, it decodes the shellcode (in place, in memory) and executes it using VirtualAlloc/RtlCopyMemory/VirtualProtect/CreateThread. The service executable is also built with customizable Windows Version Info to make it appear more like a "real" DLL.

There's probably a lot that could be cleaned up, hopefully I will get to it eventually - but for now, this works.

## Prerequisites

```
go
```

## Usage

This tool was designed to be used in Linux.

First, edit the versioninfo.template file in the goSvc directory, if desired. You can add your own description, company/copyright info., etc.

Then, from the root directory of the repository, run the generator script:

```
./generate_service_exe.sh /path/to/your/payload.bin [service_name]
```

If the optional second parameter is passed, that will be used as the name of the Windows service that gets created, as well as the name as the service executable itself. If no service name is passed, the default "youGOtserved" name will be used.

Once the service executable is created, you can drop it on a Windows machine and use the following commands:

```
# install the service that will use the executable 
# NOTE: the service bin path will be the current path of the executable when it's run
.\youGOtserved.exe install

# start the service
.\youGOtserved.exe start

# stop the service (which doesn't seem to work often when shellcode is running)
.\youGOtserved.exe stop

# Remove (delete) the service
.\youGOtserved.exe remove
```

Of course, you can install, start, stop, or delete the service using the usual Windows methods as well (sc.exe, New-Service, etc.).

## Credit

These great libraries were used in the code:
- [https://github.com/mjwhitta/babble](https://github.com/mjwhitta/babble)
- [https://github.com/josephspurrier/goversioninfo](https://github.com/josephspurrier/goversioninfo)

...and I shamelessly copied code from the [Windows service examples](https://github.com/billgraziano/go-windows-svc/tree/master) provided by the GO Project, as well as CreateThread code from this project:
- [https://github.com/Ne0nd0g/go-shellcode/tree/master/cmd/CreateThread](https://github.com/Ne0nd0g/go-shellcode/tree/master/cmd/CreateThread)
