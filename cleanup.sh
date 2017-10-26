for container in $(docker ps -a  | awk '{print $1}' ); do
  docker stop $container
  docker rm $container
done
for image in $(docker images | awk '{print $3}' ); do
  docker rmi $image
done
rm -rf dockerpanel.db
rm -rf configs
rm -rf data
