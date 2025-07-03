sudo apt install docker.io
sudo usermod -aG docker $USER
newgrp docker
docker ps