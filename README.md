# FengSheng

本项目已转移至 https://github.com/CuteReimu/TheMessage

## 声明

- **本项目采用`AGPLv3`协议开源，任何直接、间接接触本项目的软件也要求使用`AGPLv3`协议开源**
- **本项目仅用于学习和测试。不鼓励，不支持一切商业用途**
- **本项目的作者由所有参与的开发者共同所有。请尊重各位开发者的劳动成果**
- **在使用前，使用者应该对本项目有充分的了解。任何由于使用本项目提供的接口、文档等造成的不良影响和后果与任何开发者无关**
- 由于本项目的特殊性，可能随时停止开发或删档
- 本项目为开源项目，不接受任何的催单和索取行为

## 配置

```yaml
gm:
    debug_roles:  # 测试时强制设置的角色，按进入房间的顺序安排
        - 22
        - 26
    enable: false
    listen_address: 127.0.0.1:9092
listen_address: 127.0.0.1:9091  # 服务端监听端口
log:
    tcp_debug_log: true # 是否开启tcp调试日志
player:
    total_count: 5  # 玩家总人数
room_count: 200  # 最大房间数
rule:
    hand_card_count_begin: 3      # 游戏开始时摸牌数
    hand_card_count_each_turn: 3  # 每回合摸牌数
version: 1   # 需要的客户端最低版本号
```

## 游戏步骤

玩家按照逆时针的顺序依次进行回合

```mermaid
graph RL;
    id1([玩家A])-->id2([玩家B]);
    id2-->id3([玩家C]);
    id3-->id4([玩家D]);
    id4-->id5([玩家E]);
    id5-->id1;
```

每个回合只有五个阶段，且按照顺序进行

```mermaid
graph TD;
    id1([摸牌阶段])-->id2([出牌阶段]);
    id2-->id3([情报传递阶段]);
    id3-->id4([争夺阶段]);
    id4-->id5([情报接收阶段]);
```
