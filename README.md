# ctop
top for the docker container, especially for load average. "C" is for container.

The origin Linux top command only shows the information of the host machine. We want to know the load average/CPU usage of the container.

The project is origin a fork of [topic](https://github.com/silenceshell/topic), but the topic project has a serious error in calculating load average. It only calculates data of process list, not Task(Linux task_struct or thread)
