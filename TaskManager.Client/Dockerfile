# Stage 1: Build the Blazor WASM application
FROM mcr.microsoft.com/dotnet/sdk:8.0 AS build
WORKDIR /src

# Copy csproj and restore dependencies
COPY TaskManager.Client.csproj .
RUN dotnet restore

# Copy everything else and build
COPY . .
RUN dotnet publish -c Release -o /app/publish

# Stage 2: Serve with nginx
FROM nginx:alpine

# Copy custom nginx config
COPY nginx.conf /etc/nginx/nginx.conf

# Copy the published Blazor app to nginx html directory
COPY --from=build /app/publish/wwwroot /usr/share/nginx/html

# Copy appsettings for production (will be overridden by environment)
COPY wwwroot/appsettings.json /usr/share/nginx/html/appsettings.json

# Expose port 80
EXPOSE 80

# Start nginx
CMD ["nginx", "-g", "daemon off;"]
