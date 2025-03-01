FROM rabbitmq:3.13

# =========== 设置时区（ubuntu） ===========
# 设置时区，不同镜像系统之间，设置时区不同
# @link https://blog.csdn.net/wxb880114/article/details/113245119#t11
ENV TZ=Asia/Shanghai \
    DEBIAN_FRONTEND=noninteractive

RUN apt-get update \
    && apt-get install -y tzdata curl \
    && ln -fs /usr/share/zoneinfo/${TZ} /etc/localtime \
    && echo ${TZ} > /etc/timezone \
    && dpkg-reconfigure --frontend noninteractive tzdata \
    && rm -rf /var/lib/apt/lists/*
# =========== 设置时区（ubuntu） ===========

# 安装图形管理界面插件 UI port 15672
# 默认用户和密码 guest
RUN rabbitmq-plugins enable rabbitmq_management \
    && echo management_agent.disable_metrics_collector = false > /etc/rabbitmq/conf.d/management_agent.disable_metrics_collector.conf \
# 安装延迟交换机插件(不内置，须先下载)
    && curl -L -o /plugins/rabbitmq_delayed_message_exchange.ez https://github.com/rabbitmq/rabbitmq-delayed-message-exchange/releases/download/v3.13.0/rabbitmq_delayed_message_exchange-3.13.0.ez \
    && rabbitmq-plugins enable rabbitmq_delayed_message_exchange