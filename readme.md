## SuperMicro IPMI Power Supply Plugin for Zabbix Agent 2

### Contains:
- Loadable Plugin for Zabbix Agent 2
- Plugin conf file for Zabbix Agent 2
- Template for Zabbix 7.xx

### Requirements:
- Zabbix Agent 2 - version 6.0 LTS and higher
- IPMICFG Tool from SuperMicro Site for required OS: https://www.supermicro.com/en/support/resources/downloadcenter/smsdownload

### Installation 
1. Download latest release asset https://github.com/nchizhov/zabbix-supermicro-pm-plugin/releases for your OS and arch
2. Download plugin config file ```smpm.conf``` from https://github.com/nchizhov/zabbix-supermicro-pm-plugin/releases/latest/download/smpm.conf and place it to Zabbix Agent 2 ```zabbix_agent2.d/plugin.d``` folder
3. Edit ```smpm.conf```:
   - ```Plugins.SMIPMIps.System.Path``` - path to downloaded plugin executable file
   - ```Plugins.SMIPMIps.IPMITool``` - path to dowloaded IPMICFG executable file
   - ```Plugins.SMIPMIps.CollectInterval``` - IPMICFG collect data interval in minutes ```1-30```, default: ```1```
4. Additional steps for unix-like systems:
   1. Set exetubale flag for dowloaded plugin:
      ```bash 
      chmod +x plugin_file
      ```
   2. Add file to ```/etc/sudoers.d``` with name ```zabbix```:
      ```
      <zabbix_user> ALL=(ALL) NOPASSWD: <IPMICFG_path> -pminfo
      ```
      where
      - ```<zabbix_user>``` - user for Zabbix Agent 2 service\
      - ```<IPMICFG_path>``` - path to IPMICFG executable file
5. Import downloaded template ```template.xml``` from repository https://github.com/nchizhov/zabbix-supermicro-pm-plugin/releases/latest/download/template.xml to Zabbix Server Templates
6. Use imported template for needed hosts