FROM mcr.microsoft.com/dotnet/sdk:10.0-noble@sha256:de4fe90a98845d75910b9769662465c84de1a934665c431279d1382b2bddc02b
# The vulnerability audit of NuGet needs the network, and the container has none. The audit runs when the image is built, because the restore below has the network.
ENV NuGetAudit=false
ENV DOTNET_CLI_TELEMETRY_OPTOUT=1 DOTNET_NOLOGO=1 DOTNET_SKIP_FIRST_TIME_EXPERIENCE=1 NUGET_PACKAGES=/opt/nuget MSBUILDDISABLENODEREUSE=1 DOTNET_CLI_USE_MSBUILD_SERVER=0 UseSharedCompilation=false
# Restore the packages of the test project now. The container has no network, so a restore at test time can only read /opt/nuget.
COPY Restore.csproj /tmp/restore/Restore.csproj
RUN cd /tmp/restore && dotnet restore -p:NuGetAudit=true -p:NuGetAuditMode=all && rm -rf /tmp/restore /root/.nuget /root/.local && chmod -R a+rX /opt/nuget
# The powershell recipe runs the .Tests.ps1 files with Pester 5. The SDK image has pwsh. AllUsers puts the module in /usr/local/share/powershell/Modules.
RUN pwsh -NoProfile -NonInteractive -Command "Set-PSRepository PSGallery -InstallationPolicy Trusted; Install-Module Pester -RequiredVersion 5.7.1 -Scope AllUsers -Force" && rm -rf /root/.local /root/.cache
USER 1000:1000
