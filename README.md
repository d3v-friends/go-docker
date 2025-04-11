# docker engine api

* https://docs.docker.com/reference/api/engine/
* https://docs.docker.com/reference/api/engine/version/v1.48/
* 도커엔진 1.48 버전은 기존의 sdk (moby) 가 작동하지 않아서 직접 구현해야 한다.

# machine settings

- docker engine api 사용하기
- https://docs.docker.com/reference/api/engine/sdk/
- 타겟머신에서 dockerd 설정해줘야 한다.
- https://docs.docker.com/engine/daemon/remote-access/
- https://docs.docker.com/reference/api/engine/version/v1.47

~~~bash
# root 권한으로 시작한다
sudo su

# root 계정으로 변경해준다
vim /lib/systemd/system/docker.service

# 내용
# [Service] 항목에 ExecStart 파라메터에 -H tcp://0.0.0.1:2375 추가해준다
# 127.0.0.1 로 하면 로컬에서만 접속 할 수 있다.
ExecStart=/usr/bin/dockerd -H tcp://0.0.0.1:2375 -H fd:// --containerd=/run/containerd/containerd.sock $OPTIONS $DOCKER_STORAGE_OPTIONS $DOCKER_ADD_RUNTIMES

# daemon restart
systemctl daemon-reload
systemctl restart docker.service

# check port
netstat -lntp | grep dockerd

# check exec parameter
ps -ef | grep dockerd

~~~

# 기타 개발 중 필요한 정보

~~~bash
kill -9 [PID]
~~~

# docker-registry api

* https://docker-docs.uclv.cu/registry/spec/api/
 

# docker 컨테이너 내부에서 호스트 머신 접속하기
https://stackoverflow.com/questions/71668469/how-do-i-access-an-api-on-my-host-machine-from-a-docker-container

