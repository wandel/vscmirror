# Visual Studio Code Offline Mirror

Extensions are managed via the [gallery](https://github.com/microsoft/azure-devops-extension-api/tree/master/src/Gallery) feature of azure-devops-api platform

*.Icons.Default / *.Icons.Small == *.png
.VSIXPackage == *.vsix

## API

### Endpoints

#### Installers
GET https://update.code.visualstudio.com/api/update/win32-x64/stable/4949701c880d4bdb949e3c0e6b400288da7f474b
GET https://update.code.visualstudio.com/api/update/{platform}-{arch}/{channel}/{commit}

#### Extensions
GET https://main.vscode-cdn.net/extensions/marketplace.json
POST https://marketplace.visualstudio.com/_apis/public/gallery/extensionquery
GET https://marketplace.visualstudio.com/_apis/public/gallery/vscode/docker/docker/latest
GET https://marketplace.visualstudio.com/_apis/public/gallery/{publisher}/{extension}/latest

GET https://twxs.gallerycdn.vsassets.io/extensions/twxs/cmake/0.0.17/1488841920286/Microsoft.VisualStudio.Code.Manifest
GET {publisher}.gallerycdn.vsassets.io/extensions/{publisher}/{extension}/{version}/{timestamp}/{filename}

https://ms-python.gallerycdn.vsassets.io/extensions/ms-python/python/2025.4.0/1743786872189/Microsoft.VisualStudio.Code.Manifest?targetPlatform=win32-x64

POST /_apis/public/gallery/publishers/ms-vscode/extensions/cpptools/1.24.5/stats?statType=uninstall

### Types
https://github.com/microsoft/azure-devops-node-api/blob/master/api/interfaces/GalleryInterfaces.ts




## References
https://github.com/microsoft/azure-devops-extension-api/tree/master/src/Gallery




### Querys

resultMetadata

ResultCount
  TotalCount is the number of extensions that the filters matched

Categories
  entry for each category and the number of extensions that are in that category
  {"name": "Other", "count": 583},
  {"name": "Programming Languages", "count": 536},

TargetPlatforms
  entry for each target plaform and the number of extensions that are supported by that targetPlatform
  { "name": "universal", "count": 1756},
  {"name": "web", "count": 539},
