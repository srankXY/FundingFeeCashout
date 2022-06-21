#! powershell

function main {

    $checkDb=(ls | findstr ex.db | Measure-Object | findstr /i count |%{$_.split()[5]})

    if ($checkDb -ne 1){
      FundingFeeCashout-V1.1.exe init
    }

    FundingFeeCashout-V1.1.exe

    Write-Host "Please press any key to exit."
    $null = [System.Console]::ReadKey()
    break
}

main
