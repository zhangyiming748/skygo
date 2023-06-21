CREATE TABLE `asset_vehicle`
(
    `id`             int(11) NOT NULL COMMENT '主键id',
    `serial_number`  varchar(255) NOT NULL COMMENT '车型编号，程序自动生成',
    `brand`          varchar(255) NOT NULL DEFAULT '' COMMENT '车型品牌',
    `code`           varchar(255) NOT NULL DEFAULT '' COMMENT '车型代号',
    `detail`         varchar(255) NOT NULL DEFAULT '' COMMENT '车型描述',
    `create_user_id` int(11) NOT NULL COMMENT '创建用户id',
    `update_time`    int(11) unsigned NOT NULL DEFAULT '0' COMMENT '创建时间（秒）',
    `create_time`    int(11) unsigned NOT NULL DEFAULT '0' COMMENT '更新时间',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;


CREATE TABLE `asset_test_piece`
(
    `id`               int(11) NOT NULL COMMENT '主键id',
    `name`             varchar(255) NOT NULL DEFAULT '' COMMENT '测试件名称',
    `asset_vehicle_id` int(11) unsigned NOT NULL COMMENT '车型记录id',
    `detail`           varchar(255) NOT NULL DEFAULT '' COMMENT '测试件描述',
    `create_time`      int(11) unsigned NOT NULL DEFAULT '0' COMMENT '创建时间（秒）',
    `update_time`      int(11) unsigned DEFAULT '0' COMMENT '修改时间（秒）',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='资产管理-测试件';


CREATE TABLE `asset_test_piece_version`
(
    `id`                   int(11) NOT NULL COMMENT '自增主键id',
    `asset_test_piece_id`  int(11) NOT NULL COMMENT '测试件记录id',
    `version`              varchar(255)          DEFAULT NULL COMMENT '测试件版本',
    `storage_type`         tinyint(3) unsigned NOT NULL DEFAULT '1' COMMENT '存储系统分类，1mongodb',
    `create_user_id`       int(11) unsigned NOT NULL DEFAULT '0' COMMENT '创建人id',
    `update_time`          int(11) unsigned NOT NULL DEFAULT '0' COMMENT '版本记录修改时间（秒）',
    `firmware_file_uuid`   varchar(255) NOT NULL DEFAULT '' COMMENT '固件文件的唯一标识',
    `firmware_name`        varchar(255) NOT NULL DEFAULT '' COMMENT '固件名称',
    `firmware_size`        bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '固件大小（kb）',
    `firmware_device_type` tinyint(3) NOT NULL COMMENT '固件设备类型，1汽车网关(GW) 2远程通信单元(ECU) 3信息娱乐单元(IVI) ',
    `create_time`          int(11) unsigned NOT NULL DEFAULT '0' COMMENT '版本记录创建时间（秒）',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `asset_test_piece_version_file`
(
    `id`             int(11) NOT NULL COMMENT '自增主键id',
    `version_id`     varchar(255) NOT NULL COMMENT '固件某版本记录的id',
    `file_name`      varchar(255) NOT NULL COMMENT '文件名称',
    `file_size`      bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '文件大小（kb）',
    `storage_type`   tinyint(3) NOT NULL COMMENT '存储类型，1mongodb',
    `file_uuid`      varchar(255) NOT NULL DEFAULT '' COMMENT '文件存储唯一标识uuid',
    `create_time`    int(11) NOT NULL COMMENT '创建时间，即文件上传时间',
    `is_delete`      tinyint(3) unsigned NOT NULL DEFAULT '1' COMMENT '是否删除， 1否 2是',
    `delete_user_id` int(11) NOT NULL COMMENT '文件删除操作用户id',
    `delete_time`    int(11) NOT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;


CREATE TABLE `knowledge_demand`
(
    `id`             int(11) NOT NULL COMMENT '自增主键id',
    `name`           varchar(255) NOT NULL COMMENT '安全需求名称，非空',
    `category`       tinyint(3) unsigned NOT NULL COMMENT '需求类型，1企业内部标准2法规标准3渗透测试4其他',
    `code`           varchar(255) NOT NULL COMMENT '标准编号',
    `implement_time` int(11) unsigned NOT NULL COMMENT '实施日期（秒）',
    `detail`         varchar(255) NOT NULL DEFAULT '' COMMENT '描述',
    `create_time`    int(11) unsigned NOT NULL COMMENT '创建日期（秒）',
    `update_time`    int(11) unsigned NOT NULL DEFAULT '0' COMMENT '更新日期（秒）',
    `create_user_id` int(11) unsigned NOT NULL COMMENT '创建用户id',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='知识库需求';



CREATE TABLE `knowledge_demand_chapter`
(
    `id`          int(11) NOT NULL COMMENT '自增主键id',
    `code`        varchar(50)  NOT NULL COMMENT '章节编号',
    `title`       varchar(255) NOT NULL COMMENT '章节标题',
    `parent_id`   int(11) NOT NULL DEFAULT '-1' COMMENT '父章节id，-1标识无父章节',
    `parent_code` varchar(255) NOT NULL COMMENT '父章节编号',
    `content`     varchar(255) NOT NULL DEFAULT '' COMMENT '内容',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='知识库需求的章节';


CREATE TABLE `knowlege_scenario`
(
    `id`          int(11) NOT NULL COMMENT '自增主键id',
    `name`        varchar(255) NOT NULL COMMENT '安全检测场景名称',
    `demand_id`   int(11) unsigned NOT NULL COMMENT '关联知识库需求id',
    `detail`      varchar(255) NOT NULL DEFAULT '' COMMENT '场景描述',
    `create_time` int(11) unsigned NOT NULL COMMENT '创建时间',
    `tag`         varchar(255) COMMENT '标签',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='检测知识库-检测场景库';


CREATE TABLE `knowlege_scenario_chapter`
(
    `id`         int(11) NOT NULL COMMENT '自增主键id',
    `senario_id` int(11) unsigned NOT NULL COMMENT '安全检测场景id',
    `demand_id`  int(11) unsigned NOT NULL COMMENT '关联知识库需求id',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='检测知识库-检测场景库-关联条目，\r\n一个场景库关联多个需求中的条目。\r\n这样方便查询和判断条目被引用状态。';


CREATE TABLE `config_module`
(
    `id`               int(11) NOT NULL COMMENT '自增主键id',
    `module_type`      varchar(255) NOT NULL COMMENT '组件分类，如蓝牙钥匙等',
    `module_type_code` int(11) NOT NULL COMMENT '组件分类编码',
    `modele_name`      varchar(255) NOT NULL COMMENT '组件名称',
    `module_name_code` int(11) NOT NULL COMMENT '组件名称编码',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='基础配置表-组件分类';


CREATE TABLE `knowledge_test_case`
(
    `id`               int(11) NOT NULL AUTO_INCREMENT COMMENT '主键id',
    `name`             varchar(255) NOT NULL COMMENT '测试用例名称',
    `module_id`        varchar(255) NOT NULL COMMENT '测试组件/测试分类',
    `scenario_id`      int(11) NOT NULL COMMENT '安全检测场景',
    `demand_id`        int(11) NOT NULL COMMENT '关联安全需求id，需求条目存另一个表',
    `objective`        varchar(255) NOT NULL COMMENT '测试目的',
    `input`            varchar(255) NOT NULL COMMENT '外部输入',
    `test_procedure`   varchar(255) NOT NULL COMMENT '测试步骤',
    `test_standard`    varchar(255) NOT NULL COMMENT '验证标准',
    `level`            tinyint(3) NOT NULL COMMENT '测试难度',
    `test_case_level`  tinyint(3) NOT NULL COMMENT '测试用例级别',
    `test_method`      tinyint(3) NOT NULL COMMENT '测试方式，1黑盒 2灰盒 3白盒',
    `auto_test_level`  tinyint(3) NOT NULL COMMENT '自动化测试程度 1人工 2半自动化 3自动化',
    `test_tools`       varchar(255) NOT NULL COMMENT '测试工具id',
    `task_param`       varchar(255) NOT NULL COMMENT '任务参数',
    `create_user_id`   int(11) NOT NULL COMMENT '创建用户id',
    `last_op_id`       int(11) NOT NULL COMMENT '最近操作用户id',
    `create_time`      int(11) NOT NULL COMMENT '创建时间（秒）',
    `test_tool_params` varchar(255) NOT NULL COMMENT '测试任务结果',
    `content`          varchar(255) COMMENT '内容变更',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=126 DEFAULT CHARSET=utf8;

CREATE TABLE `knowledge_test_tools`
(
    `id`             int(11) NOT NULL AUTO_INCREMENT COMMENT '主键id',
    `name`           varchar(255) NOT NULL COMMENT '工具名称',
    `category`       varchar(255) NOT NULL COMMENT '工具分类',
    `introduce`      varchar(255) NOT NULL COMMENT '工具介绍',
    `brand`          varchar(255) NOT NULL COMMENT '品牌',
    `version`        varchar(255) NOT NULL COMMENT '规格型号/版本',
    `create_time`    int(11) NOT NULL COMMENT '创建时间',
    `update_time`    int(11) NOT NULL COMMENT '更新时间',
    `last_op_id`     int(11) NOT NULL COMMENT '最近操作用户id',
    `create_user_id` int(11) NOT NULL COMMENT '创建用户id',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `task`
(
    `id`                   int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键id',
    `task_uuid`            varchar(255) NOT NULL COMMENT '全局唯一id',
    `name`                 varchar(255) NOT NULL COMMENT '任务名称',
    `category`             varchar(255) NOT NULL COMMENT '类型',
    `asset_vehicle_id`     int(11) NOT NULL COMMENT '车型id',
    `piece_id`             int(11) NOT NULL COMMENT '测试件id',
    `piece_version_id`     int(11) unsigned NOT NULL COMMENT '测试件版本id',
    `scenario_id`          int(11) NOT NULL COMMENT '场景id',
    `need_connected`       int(11) NOT NULL COMMENT '是否需要连接设备， 1是 2否',
    `firmware_template_id` int(11) unsigned NOT NULL COMMENT '合规测试，测试模板id',
    `status`               int(11) NOT NULL COMMENT '状态',
    `describe`             varchar(255) NOT NULL COMMENT '描述',
    `create_user_id`       int(11) NOT NULL COMMENT '创建人id',
    `last_op_id`           int(11) NOT NULL COMMENT '最近操作用户id',
    `last_connect_time`    int(11) NOT NULL COMMENT '上次连接更新时间，单位秒',
    `hg_client_info`       varchar(255) NOT NULL COMMENT '合规上传的硬件信息，json格式',
    `hg_file_uuid`         varchar(255) NOT NULL COMMENT '合规硬件信息得到后匹配的测试用例压缩包的uuid',
    `update_time`          int(11) unsigned NOT NULL DEFAULT '0' COMMENT '更新时间',
    `create_time`          int(11) unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
    `client_info_time`     int(11) unsigned NOT NULL DEFAULT '0' COMMENT '拿到终端上传信息的时间',
    `complete_time`        int(11) unsigned NOT NULL DEFAULT '0' COMMENT '任务完成的时间',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `vulnerability`
(
    `id`             int(11) NOT NULL AUTO_INCREMENT COMMENT '主键id',
    `name`           varchar(255) NOT NULL COMMENT '漏洞名称',
    `risk_type`      tinyint(3) NOT NULL COMMENT '漏洞类型',
    `level`          tinyint(3) NOT NULL COMMENT '漏洞级别（ 0:提示 1:低危 2:中危 3:高危 4:严重 ）',
    `status`         varchar(255) NOT NULL COMMENT '漏洞状态（0:未修复 1:已修复 2:重打开）',
    `tag`            varchar(255) NOT NULL COMMENT '漏洞标签',
    `create_user_id` int(11) NOT NULL COMMENT '创建用户id',
    `last_op_id`     int(11) NOT NULL COMMENT '最近操作用户id',
    `describe`       varchar(255) NOT NULL COMMENT '漏洞描述',
    `influence`      varchar(255) NOT NULL COMMENT '影响范围',
    `sketch_map`     varchar(255) NOT NULL COMMENT '漏洞示意图',
    `suggest_id`     int(11) NOT NULL COMMENT '修复建议',
    `create_time`    int(11) NOT NULL COMMENT '创建时间',
    `update_time`    int(11) NOT NULL COMMENT '更新时间',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

--  任务的测试用例
CREATE TABLE `task_test_case`
(
    `id`                 int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键id',
    `task_id`            int(11) unsigned NOT NULL COMMENT '任务表id',
    `task_uuid`          varchar(255) NOT NULL COMMENT '全局唯一id',
    `test_case_id`       int(11) unsigned NOT NULL COMMENT '测试用例表id',
    `test_case_name`     varchar(255) NOT NULL COMMENT '测试用例名称',
    `test_result_status` tinyint(3) unsigned NOT NULL COMMENT '测试结果',
    `action_status`      tinyint(3) unsigned NOT NULL COMMENT '执行状态',
    `create_time`        int(11) unsigned NOT NULL COMMENT '创建时间',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8;

CREATE TABLE `knowledge_test_case_file`
(
    `id`             int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增主键id',
    `test_case_id`   int(11) NOT NULL COMMENT '测试用例id',
    `category`       tinyint(3) unsigned NOT NULL DEFAULT '1' COMMENT '业务类型，1测试脚本 2测试环境搭建示意图',
    `file_name`      varchar(255) NOT NULL COMMENT '文件名称',
    `file_size`      bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '文件大小（kb）',
    `storage_type`   tinyint(3) NOT NULL COMMENT '存储类型，1mongodb',
    `file_uuid`      varchar(255) NOT NULL DEFAULT '' COMMENT '文件存储唯一标识uuid',
    `create_time`    int(11) NOT NULL COMMENT '创建时间，即文件上传时间',
    `is_delete`      tinyint(3) unsigned NOT NULL DEFAULT '1' COMMENT '是否删除， 1否 2是',
    `delete_user_id` int(11) NOT NULL COMMENT '文件删除操作用户id',
    `delete_time`    int(11) NOT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `report_task`
(
    `id`            int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增主键id',
    `name`          varchar(255) COMMENT '报告名称',
    `task_id`       int(11) NOT NULL COMMENT '父类任务id',
    `status`        int(1) COMMENT '任务状态，0：未进行，1：进行中，2：成功  3：失败',
    `report_type`   int(1) COMMENT '报告类型，1：车机漏扫，2：固件扫描，3：合规检测，4：生成报告任务',
    `create_time`   varchar(255) COMMENT '创建时间 yyyymmdd',
    `complete_time` varchar(255) COMMENT '完成时间 yyyymmdd',
    `end_time`      varchar(255) COMMENT '结束时间 yyyymmdd',
    `excel_id` varbinary(255) COMMENT 'excel格式报告文件id',
    `word_id`  varbinary(255) COMMENT 'word格式报告文件id',
    `pdf_id`   varbinary(255) COMMENT 'pdf格式报告文件id',
    `excel_name` varbinary(255) COMMENT 'excel格式报告文件name',
    `word_name`  varbinary(255) COMMENT 'word格式报告文件name',
    `pdf_name`   varbinary(255) COMMENT 'pdf格式报告文件name',
    `error_msg`     varchar(512) COMMENT '失败原因',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;