Get-Content Makefile | ForEach-Object {
    if ($_ -match '^[a-zA-Z_-]+:.*?## (.*)') {
        $target = ($_.Split(':')[0]).Trim()
        $description = $matches[1].Trim()
        Write-Host ("  {0,-15} {1}" -f $target, $description)
    }
}