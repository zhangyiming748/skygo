DROP sequence IF EXISTS "public"."rule_id_seq";

CREATE SEQUENCE "public"."rule_id_seq"
INCREMENT 1
MINVALUE  1
MAXVALUE 2147483647
START 1
CACHE 1;

-- ----------------------------
-- Table structure for rule
-- ----------------------------
DROP TABLE IF EXISTS "public"."rule";
CREATE TABLE "public"."rule" (
  "id" int4 NOT NULL DEFAULT nextval('rule_id_seq'::regclass),
  "name" varchar(255) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "description" varchar(500) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "config_type" int2 NOT NULL DEFAULT 1,
  "wd_parsing_type" int2 NOT NULL DEFAULT 1,
  "wd_parsing_data" varchar(1000) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "wd_filter_content" json NOT NULL,
  "code_content" json NOT NULL,
  "category" int2 NOT NULL DEFAULT 2,
  "status" int2 NOT NULL DEFAULT 1,
  "sample_content" json NOT NULL,
  "create_time" int4
)
;
COMMENT ON COLUMN "public"."rule"."id" IS '主键id';
COMMENT ON COLUMN "public"."rule"."name" IS '解析规则名称';
COMMENT ON COLUMN "public"."rule"."description" IS '解析规则描述';
COMMENT ON COLUMN "public"."rule"."config_type" IS '配置方式，1界面，2高级指令';
COMMENT ON COLUMN "public"."rule"."wd_parsing_type" IS '界面配置，解析类型 1json';
COMMENT ON COLUMN "public"."rule"."wd_parsing_data" IS '界面配置，解析类型的辅助数据信息';
COMMENT ON COLUMN "public"."rule"."wd_filter_content" IS '界面配置。过滤规则的详细内容';
COMMENT ON COLUMN "public"."rule"."code_content" IS '指令配置，整个配置内容';
COMMENT ON COLUMN "public"."rule"."category" IS '分类：1预定义，2自定义';
COMMENT ON COLUMN "public"."rule"."status" IS '状态，1草稿  2正常';
COMMENT ON COLUMN "public"."rule"."sample_content" IS '样本列表，数据是json数组';
COMMENT ON COLUMN "public"."rule"."create_time" IS '创建时间，秒';
COMMENT ON TABLE "public"."rule" IS '解析规则';

-- ----------------------------
-- Records of rule
-- ----------------------------
INSERT INTO "public"."rule" VALUES (5, '流量-会话关闭上报', '流量-会话关闭上报', 1, 1, '1', '[{
	"ip_private_parse-parse": {
		"field": "sip",
		"hatch_field": "sip_is_lan",
		"is_private_value": "1",
		"not_private_value": "0"
	}
}, {
	"ip_private_parse-parse": {
		"field": "dip",
		"hatch_field": "dip_is_lan",
		"is_private_value": "1",
		"not_private_value": "0"
	}
}, {
	"call_func-add_field_by_func": {
		"field": "receive_time",
		"func_name": "millisecond_timestamp"
	}
}, {
	"call_func-add_field_by_func": {
		"field": "log_id",
		"func_name": "wujiang_log_id"
	}
}, {
	"api_ip_hatch": {
		"field": "sip",
		"set_fields": {
			"country_cn": "geo_sip_country_cn",
			"country_en": "geo_sip_country_en",
			"province": "geo_sip_province",
			"city": "geo_sip_city",
			"districts": "geo_sip_district",
			"latitude": "geo_sip_latitude",
			"longitude": "geo_sip_longitude"
		}
	}
}, {
	"api_ip_hatch": {
		"field": "dip",
		"set_fields": {
			"country_cn": "geo_dip_country_cn",
			"country_en": "geo_dip_country_en",
			"province": "geo_dip_province",
			"city": "geo_dip_city",
			"districts": "geo_dip_district",
			"latitude": "geo_dip_latitude",
			"longitude": "geo_dip_longitude"
		}
	}
}, {
	"call_func-add_field_by_func_param": {
		"field": "md5",
		"func_name": "md5",
		"func_param": "%{sip}%{dip}%{protocol}%{s_port}%{d_port}"
	}
}, {
	"mutate-del_field": [{
		"field": "etl_origin_raw_log"
	}]
}]', '[]', 2, 2, '[]', 0);
INSERT INTO "public"."rule" VALUES (8, '威胁数据上报', '威胁数据上报', 1, 1, '1', '[ {
	"mutate-define_field": [{
		"field": "sample_type",
		"value": ""
	}, {
		"field": "sample_md5",
		"value": ""
	}, {
		"field": "sample_sha256",
		"value": ""
	}, {
		"field": "sample_name",
		"value": ""
	}, {
		"field": "sample_direction",
		"value": ""
	}, {
		"field": "sample_size",
		"value": "0"
	}, {
		"field": "is_ioc",
		"value": "0"
	}, {
		"field": "ioc",
		"value": ""
	}, {
		"field": "ioc_detail",
		"value": ""
	}, {
		"field": "ioc_type",
		"value": ""
	}, {
		"field": "ioc_source",
		"value": ""
	}, {
		"field": "ioc_malicious_family",
		"value": ""
	}, {
		"field": "ioc_campaign",
		"value": ""
	}, {
		"field": "ioc_control_type",
		"value": ""
	}, {
		"field": "is_sample",
		"value": "0"
	}]
}, {
	"ip_private_parse-parse": {
		"field": "sip",
		"hatch_field": "sip_is_lan",
		"is_private_value": "1",
		"not_private_value": "0"
	}
}, {
	"ip_private_parse-parse": {
		"field": "dip",
		"hatch_field": "dip_is_lan",
		"is_private_value": "1",
		"not_private_value": "0"
	}
}, {
	"call_func-add_field_by_func": {
		"field": "receive_time",
		"func_name": "millisecond_timestamp"
	}
}, {
	"call_func-add_field_by_func": {
		"field": "log_id",
		"func_name": "wujiang_log_id"
	}
}, {
	"api_ip_hatch": {
		"field": "sip",
		"set_fields": {
			"country_cn": "geo_sip_country_cn",
			"country_en": "geo_sip_country_en",
			"province": "geo_sip_province",
			"city": "geo_sip_city",
			"districts": "geo_sip_district",
			"latitude": "geo_sip_latitude",
			"longitude": "geo_sip_longitude"
		}
	}
}, {
	"api_ip_hatch": {
		"field": "dip",
		"set_fields": {
			"country_cn": "geo_dip_country_cn",
			"country_en": "geo_dip_country_en",
			"province": "geo_dip_province",
			"city": "geo_dip_city",
			"districts": "geo_dip_district",
			"latitude": "geo_dip_latitude",
			"longitude": "geo_dip_longitude"
		}
	}
},{
	"mutate-del_field":[
		{"field":"etl_origin_raw_log"}
	]
}]', '[]', 2, 2, '[]', 0);
INSERT INTO "public"."rule" VALUES (11, '中能融合-警报日志上报', '中能融合-警报日志上报', 1, 1, '1', '[{
	"mutate-define_field": [{
		"field": "is_ipv6",
		"value": "0"
	}, {
		"field": "protocol",
		"value": ""
	}, {
		"field": "service",
		"value": ""
	}, {
		"field": "app_name",
		"value": ""
	}, {
		"field": "attacker_ip",
		"value": "%{src_ip}"
	}, {
		"field": "victim_ip",
		"value": "%{dest_ip}"
	}, {
		"field": "solution",
		"value": ""
	}, {
		"field": "rule_id",
		"value": ""
	}, {
		"field": "attack_stage",
		"value": "未知"
	}, {
		"field": "attack_result",
		"value": "未知"
	}, {
		"field": "attack_confidence",
		"value": "中"
	}, {
		"field": "is_sample",
		"value": "0"
	}, {
		"field": "sample_type",
		"value": ""
	}, {
		"field": "sample_md5",
		"value": ""
	}, {
		"field": "sample_sha256",
		"value": ""
	}, {
		"field": "sample_name",
		"value": ""
	}, {
		"field": "sample_direction",
		"value": ""
	}, {
		"field": "sample_size",
		"value": "0"
	}, {
		"field": "is_ioc",
		"value": "0"
	}, {
		"field": "ioc",
		"value": ""
	}, {
		"field": "ioc_detail",
		"value": ""
	}, {
		"field": "ioc_type",
		"value": ""
	}, {
		"field": "ioc_source",
		"value": ""
	}, {
		"field": "ioc_malicious_family",
		"value": ""
	}, {
		"field": "ioc_campaign",
		"value": ""
	}, {
		"field": "ioc_control_type",
		"value": ""
	}, {
		"field": "compromise_state",
		"value": "未知"
	}, {
		"field": "data_source",
		"value": "zhongneng"
	}, {
		"field": "log_type",
		"value": "Threat"
	}]
}, {
	"mutate-merge_field": {
		"merge_fields": ["device_id", "netlink_id"],
		"merge_no": "",
		"new_field": "serial_num"
	}
}, {
	"mutate-rename":[{
		"field": "sip",
		"value": "src_ip"
	}, {
		"field": "dip",
		"value": "dest_ip"
	}, {
		"field": "s_port",
		"value": "src_port"
	}, {
		"field": "d_port",
		"value": "dest_port"
	}, {
		"field": "found_time",
		"value": "stat_time"
	}, {
		"field": "threat_type",
		"value": "alarm_category"
	}, {
		"field": "threat_desc",
		"value": "alarm_description"
	}, {
		"field": "threat_name",
		"value": "alarm_name"
	}]
}, {
	"ip_private_parse-parse": {
		"field": "sip",
		"hatch_field": "sip_is_lan",
		"is_private_value": "1",
		"not_private_value": "0"
	}
}, {
	"ip_private_parse-parse": {
		"field": "dip",
		"hatch_field": "dip_is_lan",
		"is_private_value": "1",
		"not_private_value": "0"
	}
}, {
	"call_func-add_field_by_func": {
		"field": "receive_time",
		"func_name": "millisecond_timestamp"
	}
}, {
	"call_func-add_field_by_func": {
		"field": "log_id",
		"func_name": "wujiang_log_id"
	}
}, {
	"mutate-field_map_rel": {
		"field": "alarm_severity",
		"new_field": "hazard_level",
		"map_rel": [{
			"value": "1",
			"new_value": "低危"
		}, {
			"value": "2",
			"new_value": "中危"
		}, {
			"value": "3",
			"new_value": "高危"
		}]
	}
}, {
	"mutate-del_field": [{
		"field": "etl_origin_raw_log"
	}]
}]', '[]', 2, 2, '[]', 0);
INSERT INTO "public"."rule" VALUES (10, '中能融合-行为日志上报', '中能融合-行为日志上报', 1, 1, '1', '[{
	"call_func-add_field_by_func": {
		"field": "receive_time",
		"func_name": "millisecond_timestamp"
	}
},{
	"mutate-del_field":[
		{"field":"etl_origin_raw_log"}
	]
}]', '[]', 2, 2, '[]', 0);
INSERT INTO "public"."rule" VALUES (9, '资产上报', '资产上报', 1, 1, '1', '[{
	"ip_private_parse-parse": {
		"field": "ip",
		"hatch_field": "ip_is_lan",
		"is_private_value": "1",
		"not_private_value": "0"
	}
}, {
	"call_func-add_field_by_func": {
		"field": "log_id",
		"func_name": "wujiang_log_id"
	}
}, {
	"api_ip_hatch": {
		"field": "ip",
		"set_fields": {
			"country_cn": "geo_ip_country_cn",
			"country_en": "geo_ip_country_en",
			"province": "geo_ip_province",
			"city": "geo_ip_city",
			"districts": "geo_ip_district",
			"latitude": "geo_ip_latitude",
			"longitude": "geo_ip_longitude"
		}
	}
}, {
	"mutate-json_decode_add_field": [{
		"field": "asset_vul",
		"value": "[]"
	}]

},{
	"mutate-del_field":[
		{"field":"etl_origin_raw_log"}
	]
}]', '[]', 2, 2, '["{\"active_time\":1600771725000,\"asset_sn\":\"\",\"category\":\"\",\"data_source\":\"ITD\",\"is_ipv6\":\"0\",\"log_type\":\"Asset\",\"mac\":\"\",\"name\":\"\",\"open_port\":[{\"app_name\":\"https\",\"port\":443,\"protocol\":\"TCP\",\"service\":\"https\"},{\"app_name\":\"ftp\",\"port\":21,\"protocol\":\"TCP\",\"service\":\"ftp\"}],\"os\":\"Windows\",\"os_version\":\"XP\",\"product\":\"\",\"serial_num\":\"itd_device_serial_001\",\"type\":\"\",\"vendor\":\"\",\"version\":\"\",\"ip\":\"104.192.108.9\",\"found_time\":1615279195108}"]', 0);
INSERT INTO "public"."rule" VALUES (1, '域名解析上报', '域名解析上报', 1, 1, '1', '[{
	"ip_private_parse-parse": {
		"field": "sip",
		"hatch_field": "sip_is_lan",
		"is_private_value": "1",
		"not_private_value": "0"
	}
}, {
	"ip_private_parse-parse": {
		"field": "dip",
		"hatch_field": "dip_is_lan",
		"is_private_value": "1",
		"not_private_value": "0"
	}
}, {
	"call_func-add_field_by_func": {
		"field": "receive_time",
		"func_name": "millisecond_timestamp"
	}
}, {
	"api_ip_hatch": {
		"field": "sip",
		"set_fields": {
			"country_cn": "geo_sip_country_cn",
			"country_en": "geo_sip_country_en",
			"province": "geo_sip_province",
			"city": "geo_sip_city",
			"districts": "geo_sip_district",
			"latitude": "geo_sip_latitude",
			"longitude": "geo_sip_longitude"
		}
	}
}, {
	"api_ip_hatch": {
		"field": "dip",
		"set_fields": {
			"country_cn": "geo_dip_country_cn",
			"country_en": "geo_dip_country_en",
			"province": "geo_dip_province",
			"city": "geo_dip_city",
			"districts": "geo_dip_district",
			"latitude": "geo_dip_latitude",
			"longitude": "geo_dip_longitude"
		}
	}
},{
	"mutate-del_field":[
		{"field":"etl_origin_raw_log"}
	]
}]', '[]', 2, 2, '[]', 0);
INSERT INTO "public"."rule" VALUES (7, '工业关键事件上报', '工业关键事件上报', 1, 1, '1', '[{
	"ip_private_parse-parse": {
		"field": "sip",
		"hatch_field": "sip_is_lan",
		"is_private_value": "1",
		"not_private_value": "0"
	}
}, {
	"ip_private_parse-parse": {
		"field": "dip",
		"hatch_field": "dip_is_lan",
		"is_private_value": "1",
		"not_private_value": "0"
	}
}, {
	"call_func-add_field_by_func": {
		"field": "receive_time",
		"func_name": "millisecond_timestamp"
	}
}, {
	"call_func-add_field_by_func": {
		"field": "log_id",
		"func_name": "wujiang_log_id"
	}
}, {
	"api_ip_hatch": {
		"field": "sip",
		"set_fields": {
			"country_cn": "geo_sip_country_cn",
			"country_en": "geo_sip_country_en",
			"province": "geo_sip_province",
			"city": "geo_sip_city",
			"districts": "geo_sip_district",
			"latitude": "geo_sip_latitude",
			"longitude": "geo_sip_longitude"
		}
	}
}, {
	"api_ip_hatch": {
		"field": "dip",
		"set_fields": {
			"country_cn": "geo_dip_country_cn",
			"country_en": "geo_dip_country_en",
			"province": "geo_dip_province",
			"city": "geo_dip_city",
			"districts": "geo_dip_district",
			"latitude": "geo_dip_latitude",
			"longitude": "geo_dip_longitude"
		}
	}
},{
	"mutate-del_field":[
		{"field":"etl_origin_raw_log"}
	]
}]', '[]', 2, 2, '[]', 0);
INSERT INTO "public"."rule" VALUES (3, 'WEB访问上报', 'WEB访问上报', 1, 1, '1', '[{
		"ip_private_parse-parse": {
			"field": "sip",
			"hatch_field": "sip_is_lan",
			"is_private_value": "1",
			"not_private_value": "0"
		}
	}, {
		"ip_private_parse-parse": {
			"field": "dip",
			"hatch_field": "dip_is_lan",
			"is_private_value": "1",
			"not_private_value": "0"
		}
	}, {
		"call_func-add_field_by_func": {
			"field": "receive_time",
			"func_name": "millisecond_timestamp"
		}
	}, {
		"call_func-add_field_by_func": {
			"field": "log_id",
			"func_name": "wujiang_log_id"
		}
	}, {
		"api_ip_hatch": {
			"field": "sip",
			"set_fields": {
				"country_cn": "geo_sip_country_cn",
				"country_en": "geo_sip_country_en",
				"province": "geo_sip_province",
				"city": "geo_sip_city",
				"districts": "geo_sip_district",
				"latitude": "geo_sip_latitude",
				"longitude": "geo_sip_longitude"
			}
		}
	}, {
		"api_ip_hatch": {
			"field": "dip",
			"set_fields": {
				"country_cn": "geo_dip_country_cn",
				"country_en": "geo_dip_country_en",
				"province": "geo_dip_province",
				"city": "geo_dip_city",
				"districts": "geo_dip_district",
				"latitude": "geo_dip_latitude",
				"longitude": "geo_dip_longitude"
			}
		}
	},{
	"mutate-del_field":[
		{"field":"etl_origin_raw_log"}
	]
}]', '[]', 2, 2, '[]', 0);
INSERT INTO "public"."rule" VALUES (4, '白名单数据上报', '白名单数据上报', 1, 1, '1', '[{
	"ip_private_parse-parse": {
		"field": "sip",
		"hatch_field": "sip_is_lan",
		"is_private_value": "1",
		"not_private_value": "0"
	}
}, {
	"ip_private_parse-parse": {
		"field": "dip",
		"hatch_field": "dip_is_lan",
		"is_private_value": "1",
		"not_private_value": "0"
	}
}, {
	"call_func-add_field_by_func": {
		"field": "log_id",
		"func_name": "wujiang_log_id"
	}
}, {
	"call_func-add_field_by_func": {
		"field": "receive_time",
		"func_name": "millisecond_timestamp"
	}
}, {
	"api_ip_hatch": {
		"field": "sip",
		"set_fields": {
			"country_cn": "geo_sip_country_cn",
			"country_en": "geo_sip_country_en",
			"province": "geo_sip_province",
			"city": "geo_sip_city",
			"districts": "geo_sip_district",
			"latitude": "geo_sip_latitude",
			"longitude": "geo_sip_longitude"
		}
	}
}, {
	"api_ip_hatch": {
		"field": "dip",
		"set_fields": {
			"country_cn": "geo_dip_country_cn",
			"country_en": "geo_dip_country_en",
			"province": "geo_dip_province",
			"city": "geo_dip_city",
			"districts": "geo_dip_district",
			"latitude": "geo_dip_latitude",
			"longitude": "geo_dip_longitude"
		}
	}
},{
	"mutate-del_field":[
		{"field":"etl_origin_raw_log"}
	]
}]', '[]', 2, 2, '[]', 0);
INSERT INTO "public"."rule" VALUES (6, '邮件行为上报', '邮件行为上报', 1, 1, '1', '[{
	"ip_private_parse-parse": {
		"field": "sip",
		"hatch_field": "sip_is_lan",
		"is_private_value": "1",
		"not_private_value": "0"
	}
}, {
	"ip_private_parse-parse": {
		"field": "dip",
		"hatch_field": "dip_is_lan",
		"is_private_value": "1",
		"not_private_value": "0"
	}
}, {
	"call_func-add_field_by_func": {
		"field": "receive_time",
		"func_name": "millisecond_timestamp"
	}
}, {
	"api_ip_hatch": {
		"field": "sip",
		"set_fields": {
			"country_cn": "geo_sip_country_cn",
			"country_en": "geo_sip_country_en",
			"province": "geo_sip_province",
			"city": "geo_sip_city",
			"districts": "geo_sip_district",
			"latitude": "geo_sip_latitude",
			"longitude": "geo_sip_longitude"
		}
	}
}, {
	"api_ip_hatch": {
		"field": "dip",
		"set_fields": {
			"country_cn": "geo_dip_country_cn",
			"country_en": "geo_dip_country_en",
			"province": "geo_dip_province",
			"city": "geo_dip_city",
			"districts": "geo_dip_district",
			"latitude": "geo_dip_latitude",
			"longitude": "geo_dip_longitude"
		}
	}
},{
	"mutate-del_field":[
		{"field":"etl_origin_raw_log"}
	]
}]', '[]', 2, 2, '[]', 0);
INSERT INTO "public"."rule" VALUES (2, 'SSL加密协商上报', 'SSL加密协商上报', 1, 1, '1', '[{
	"ip_private_parse-parse": {
		"field": "sip",
		"hatch_field": "sip_is_lan",
		"is_private_value": "1",
		"not_private_value": "0"
	}
}, {
	"ip_private_parse-parse": {
		"field": "dip",
		"hatch_field": "dip_is_lan",
		"is_private_value": "1",
		"not_private_value": "0"
	}
}, {
	"call_func-add_field_by_func": {
		"field": "receive_time",
		"func_name": "millisecond_timestamp"
	}
}, {
	"call_func-add_field_by_func": {
		"field": "log_id",
		"func_name": "wujiang_log_id"
	}
}, {
	"api_ip_hatch": {
		"field": "sip",
		"set_fields": {
			"country_cn": "geo_sip_country_cn",
			"country_en": "geo_sip_country_en",
			"province": "geo_sip_province",
			"city": "geo_sip_city",
			"districts": "geo_sip_district",
			"latitude": "geo_sip_latitude",
			"longitude": "geo_sip_longitude"
		}
	}
}, {
	"api_ip_hatch": {
		"field": "dip",
		"set_fields": {
			"country_cn": "geo_dip_country_cn",
			"country_en": "geo_dip_country_en",
			"province": "geo_dip_province",
			"city": "geo_dip_city",
			"districts": "geo_dip_district",
			"latitude": "geo_dip_latitude",
			"longitude": "geo_dip_longitude"
		}
	}
},{
	"mutate-del_field":[
		{"field":"etl_origin_raw_log"}
	]
}]', '[]', 2, 2, '[]', 0);

-- ----------------------------
-- Primary Key structure for table rule
-- ----------------------------
ALTER TABLE "public"."rule" ADD CONSTRAINT "_copy_4" PRIMARY KEY ("id");
