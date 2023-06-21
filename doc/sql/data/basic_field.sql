DROP sequence IF EXISTS "public"."basic_field_id_seq";

CREATE SEQUENCE "public"."basic_field_id_seq"
INCREMENT 1
MINVALUE  1
MAXVALUE 2147483647
START 1
CACHE 1;

SELECT setval('"public"."basic_field_id_seq"', 690, true);

-- ----------------------------
-- Table structure for basic_field
-- ----------------------------
DROP TABLE IF EXISTS "public"."basic_field";
CREATE TABLE "public"."basic_field" (
  "id" int4 NOT NULL DEFAULT nextval('basic_field_id_seq'::regclass),
  "name" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "alias" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "type" int2 NOT NULL,
  "description" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "category" int2 NOT NULL,
  "create_time" int4
)
;
COMMENT ON COLUMN "public"."basic_field"."id" IS '自增主键';
COMMENT ON COLUMN "public"."basic_field"."name" IS '字段名称';
COMMENT ON COLUMN "public"."basic_field"."alias" IS '字段别名,即中文名';
COMMENT ON COLUMN "public"."basic_field"."type" IS '字段类型';
COMMENT ON COLUMN "public"."basic_field"."description" IS '描述';
COMMENT ON COLUMN "public"."basic_field"."category" IS '分类：1预定义，2自定义';
COMMENT ON COLUMN "public"."basic_field"."create_time" IS '创建时间';
COMMENT ON TABLE "public"."basic_field" IS '基础字段';

-- ----------------------------
-- Records of basic_field
-- ----------------------------
INSERT INTO "public"."basic_field" VALUES (566, 'bytes_sent', 'bytes_sent', 4, '', 1, 1614684337);
INSERT INTO "public"."basic_field" VALUES (567, 'xff', 'xff', 2, '', 1, 1614684337);
INSERT INTO "public"."basic_field" VALUES (568, 'attack_result', 'attack_result', 1, '', 1, 1614684337);
INSERT INTO "public"."basic_field" VALUES (569, 'sample_md5', 'sample_md5', 2, '', 1, 1614684337);
INSERT INTO "public"."basic_field" VALUES (570, 'db_name', 'db_name', 2, '', 1, 1614684337);
INSERT INTO "public"."basic_field" VALUES (571, 'attack_confidence', 'attack_confidence', 1, '', 1, 1614684337);
INSERT INTO "public"."basic_field" VALUES (572, 'cname', 'cname', 2, '', 1, 1614684337);
INSERT INTO "public"."basic_field" VALUES (573, 'app_name', 'app_name', 1, '', 1, 1614684337);
INSERT INTO "public"."basic_field" VALUES (574, 'dip_is_lan', 'dip_is_lan', 1, '', 1, 1614684337);
INSERT INTO "public"."basic_field" VALUES (575, 'from', 'from', 2, '', 1, 1614684337);
INSERT INTO "public"."basic_field" VALUES (576, 'attach_md5', 'attach_md5', 1, '', 1, 1614684337);
INSERT INTO "public"."basic_field" VALUES (577, 'log_type', 'log_type', 1, '', 1, 1614684337);
INSERT INTO "public"."basic_field" VALUES (578, 'attack_stage', 'attack_stage', 1, '', 1, 1614684337);
INSERT INTO "public"."basic_field" VALUES (579, 'service', 'service', 1, '', 1, 1614684337);
INSERT INTO "public"."basic_field" VALUES (580, 'file_name', 'file_name', 1, '', 1, 1614684337);
INSERT INTO "public"."basic_field" VALUES (581, 'file_type', 'file_type', 2, '', 1, 1614684337);
INSERT INTO "public"."basic_field" VALUES (582, 'pkts_send', 'pkts_send', 4, '', 1, 1614684337);
INSERT INTO "public"."basic_field" VALUES (583, 'url', 'url', 2, '', 1, 1614684337);
INSERT INTO "public"."basic_field" VALUES (584, 'threat_class', 'threat_class', 1, '', 1, 1614684337);
INSERT INTO "public"."basic_field" VALUES (585, 'ioc_control_type', 'ioc_control_type', 2, '', 1, 1614684337);
INSERT INTO "public"."basic_field" VALUES (586, 'geo_sip_province', 'geo_sip_province', 1, '', 1, 1614684337);
INSERT INTO "public"."basic_field" VALUES (587, 'd_port', 'd_port', 4, '', 1, 1614684337);
INSERT INTO "public"."basic_field" VALUES (588, 'protocol', 'protocol', 1, '', 1, 1614684337);
INSERT INTO "public"."basic_field" VALUES (589, 'ics_event', 'ics_event', 2, '', 1, 1614684337);
INSERT INTO "public"."basic_field" VALUES (590, 'threat_desc', 'threat_desc', 2, '', 1, 1614684337);
INSERT INTO "public"."basic_field" VALUES (591, 'ioc_detail', 'ioc_detail', 2, '', 1, 1614684337);
INSERT INTO "public"."basic_field" VALUES (592, 'reply_code', 'reply_code', 1, '', 1, 1614684337);
INSERT INTO "public"."basic_field" VALUES (593, 'attachment', 'attachment', 10, '', 1, 1614684337);
INSERT INTO "public"."basic_field" VALUES (594, 'is_ioc', 'is_ioc', 1, '', 1, 1614684337);
INSERT INTO "public"."basic_field" VALUES (602, 'ioc', 'ioc', 2, '', 1, 1614684337);
INSERT INTO "public"."basic_field" VALUES (681, 'stat_time', '分析时间', 8, '', 1, 1614684337);
INSERT INTO "public"."basic_field" VALUES (682, 'behavior_level', '警报等级', 1, '', 1, 1614684337);
INSERT INTO "public"."basic_field" VALUES (683, 'behavior_src_ip', '源IP', 1, '', 1, 1614684337);
INSERT INTO "public"."basic_field" VALUES (684, 'behavior_dest_ip', '目标IP', 1, '', 1, 1614684337);
INSERT INTO "public"."basic_field" VALUES (685, 'behavior_category', '行为分类', 1, '', 1, 1614684337);
INSERT INTO "public"."basic_field" VALUES (686, 'behavior_data', '行为数据', 2, '', 1, 1614684337);
INSERT INTO "public"."basic_field" VALUES (687, 'behavior_netobj_count', '网络对象个数', 5, '', 1, 1614684337);
INSERT INTO "public"."basic_field" VALUES (688, 'behavior_type', '行为类型', 1, '', 1, 1614684337);
INSERT INTO "public"."basic_field" VALUES (689, 'behavior_name', '行为名称', 1, '', 1, 1614684337);
INSERT INTO "public"."basic_field" VALUES (690, 'behavior_desc', '行为描述', 2, '', 1, 1614684337);
INSERT INTO "public"."basic_field" VALUES (595, 'dispose_suggest', 'dispose_suggest', 2, '', 1, 1614684337);
INSERT INTO "public"."basic_field" VALUES (596, 's_port', 's_port', 4, '', 1, 1614684337);
INSERT INTO "public"."basic_field" VALUES (597, 'passwd', 'passwd', 2, '', 1, 1614684337);
INSERT INTO "public"."basic_field" VALUES (598, 'referer', 'referer', 2, '', 1, 1614684337);
INSERT INTO "public"."basic_field" VALUES (599, 'geo_sip_city', 'geo_sip_city', 1, '', 1, 1614684337);
INSERT INTO "public"."basic_field" VALUES (600, 'detect_method', 'detect_method', 1, '', 1, 1614684337);
INSERT INTO "public"."basic_field" VALUES (601, 'file_sha256', 'file_sha256', 2, '', 1, 1614684337);
INSERT INTO "public"."basic_field" VALUES (603, 'addr', 'addr', 2, '', 1, 1614684337);
INSERT INTO "public"."basic_field" VALUES (604, 'geo_dip_latitude', 'geo_dip_latitude', 6, '', 1, 1614684337);
INSERT INTO "public"."basic_field" VALUES (605, 'subject', 'subject', 2, '', 1, 1614684337);
INSERT INTO "public"."basic_field" VALUES (606, 'geo_dip_district', 'geo_dip_district', 1, '', 1, 1614684337);
INSERT INTO "public"."basic_field" VALUES (607, 'threat_name', 'threat_name', 2, '', 1, 1614684337);
INSERT INTO "public"."basic_field" VALUES (608, 'attacker_ip', 'attacker_ip', 9, '', 1, 1614684337);
INSERT INTO "public"."basic_field" VALUES (609, 'mx', 'mx', 2, '', 1, 1614684337);
INSERT INTO "public"."basic_field" VALUES (610, 'server_name', 'server_name', 2, '', 1, 1614684337);
INSERT INTO "public"."basic_field" VALUES (611, 'notafter', 'notafter', 1, '', 1, 1614684337);
INSERT INTO "public"."basic_field" VALUES (612, 'origin', 'origin', 2, '', 1, 1614684337);
INSERT INTO "public"."basic_field" VALUES (613, 'geo_sip_district', 'geo_sip_district', 1, '', 1, 1614684337);
INSERT INTO "public"."basic_field" VALUES (614, 'x-forward', 'x-forward', 2, '', 1, 1614684337);
INSERT INTO "public"."basic_field" VALUES (615, 'direction', 'direction', 1, '', 1, 1614684337);
INSERT INTO "public"."basic_field" VALUES (616, 'sql_info', 'sql_info', 2, '', 1, 1614684337);
INSERT INTO "public"."basic_field" VALUES (617, 'content-type', 'content-type', 1, '', 1, 1614684337);
INSERT INTO "public"."basic_field" VALUES (618, 'dip', 'dip', 9, '', 1, 1614684337);
INSERT INTO "public"."basic_field" VALUES (619, 'ioc_source', 'ioc_source', 2, '', 1, 1614684337);
INSERT INTO "public"."basic_field" VALUES (620, 'user_name', 'user_name', 2, '', 1, 1614684337);
INSERT INTO "public"."basic_field" VALUES (621, 'sample_type', 'sample_type', 1, '', 1, 1614684337);
INSERT INTO "public"."basic_field" VALUES (622, 'file_key', 'file_key', 2, '', 1, 1614684337);
INSERT INTO "public"."basic_field" VALUES (623, 'bytes_received', 'bytes_received', 4, '', 1, 1614684337);
INSERT INTO "public"."basic_field" VALUES (624, 'dns_type', 'dns_type', 1, '', 1, 1614684337);
INSERT INTO "public"."basic_field" VALUES (625, 'victim_ip', 'victim_ip', 9, '', 1, 1614684337);
INSERT INTO "public"."basic_field" VALUES (626, 'sip', 'sip', 9, '', 1, 1614684337);
INSERT INTO "public"."basic_field" VALUES (627, 'sample_name', 'sample_name', 2, '', 1, 1614684337);
INSERT INTO "public"."basic_field" VALUES (628, 'public_key', 'public_key', 1, '', 1, 1614684337);
INSERT INTO "public"."basic_field" VALUES (629, 'geo_sip_latitude', 'geo_sip_latitude', 6, '', 1, 1614684337);
INSERT INTO "public"."basic_field" VALUES (630, 'ioc_campaign', 'ioc_campaign', 2, '', 1, 1614684337);
INSERT INTO "public"."basic_field" VALUES (631, 'ioc_type', 'ioc_type', 2, '', 1, 1614684337);
INSERT INTO "public"."basic_field" VALUES (632, 'duration', 'duration', 4, '', 1, 1614684337);
INSERT INTO "public"."basic_field" VALUES (633, 'data_source', 'data_source', 2, '', 1, 1614684337);
INSERT INTO "public"."basic_field" VALUES (634, 'serial_num', 'serial_num', 2, '', 1, 1614684337);
INSERT INTO "public"."basic_field" VALUES (636, 'offence_value', 'offence_value', 2, '', 1, 1614684337);
INSERT INTO "public"."basic_field" VALUES (637, 'is_ipv6', 'is_ipv6', 1, '', 1, 1614684337);
INSERT INTO "public"."basic_field" VALUES (639, 'geo_sip_country', 'geo_sip_country', 1, '', 1, 1614684337);
INSERT INTO "public"."basic_field" VALUES (640, 'hazard_level', 'hazard_level', 1, '', 1, 1614684337);
INSERT INTO "public"."basic_field" VALUES (641, 'geo_dip_city', 'geo_dip_city', 1, '', 1, 1614684337);
INSERT INTO "public"."basic_field" VALUES (642, 'event_class', 'event_class', 1, '', 1, 1614684337);
INSERT INTO "public"."basic_field" VALUES (643, 'geo_dip_province', 'geo_dip_province', 1, '', 1, 1614684337);
INSERT INTO "public"."basic_field" VALUES (644, 'compromise_state', 'compromise_state', 1, '', 1, 1614684337);
INSERT INTO "public"."basic_field" VALUES (645, 'sample_size', 'sample_size', 4, '', 1, 1614684337);
INSERT INTO "public"."basic_field" VALUES (646, 'cert', 'cert', 10, '', 1, 1614684337);
INSERT INTO "public"."basic_field" VALUES (647, 'ioc_malicious_family', 'ioc_malicious_family', 2, '', 1, 1614684337);
INSERT INTO "public"."basic_field" VALUES (648, 'sample_sha256', 'sample_sha256', 2, '', 1, 1614684337);
INSERT INTO "public"."basic_field" VALUES (649, 'file_size', 'file_size', 4, '', 1, 1614684337);
INSERT INTO "public"."basic_field" VALUES (650, 'cc', 'cc', 2, '', 1, 1614684337);
INSERT INTO "public"."basic_field" VALUES (651, 'name', 'name', 1, '', 1, 1614684337);
INSERT INTO "public"."basic_field" VALUES (652, 'notbefore', 'notbefore', 1, '', 1, 1614684337);
INSERT INTO "public"."basic_field" VALUES (653, 'geo_sip_longitude', 'geo_sip_longitude', 6, '', 1, 1614684337);
INSERT INTO "public"."basic_field" VALUES (654, 'user-agent', 'user-agent', 2, '', 1, 1614684337);
INSERT INTO "public"."basic_field" VALUES (655, 'offence_type', 'offence_type', 2, '', 1, 1614684337);
INSERT INTO "public"."basic_field" VALUES (656, 'cookie', 'cookie', 2, '', 1, 1614684337);
INSERT INTO "public"."basic_field" VALUES (657, 'solution', 'solution', 1, '', 1, 1614684337);
INSERT INTO "public"."basic_field" VALUES (658, 'is_sample', 'is_sample', 1, '', 1, 1614684337);
INSERT INTO "public"."basic_field" VALUES (659, 'geo_dip_longitude', 'geo_dip_longitude', 6, '', 1, 1614684337);
INSERT INTO "public"."basic_field" VALUES (661, 'sip_is_lan', 'sip_is_lan', 1, '', 1, 1614684337);
INSERT INTO "public"."basic_field" VALUES (662, 'pkts_received', 'pkts_received', 4, '', 1, 1614684337);
INSERT INTO "public"."basic_field" VALUES (663, 'ics_protocol', 'ics_protocol', 2, '', 1, 1614684337);
INSERT INTO "public"."basic_field" VALUES (664, 'rule_id', 'rule_id', 1, '', 1, 1614684337);
INSERT INTO "public"."basic_field" VALUES (665, 'file_md5', 'file_md5', 2, '', 1, 1614684337);
INSERT INTO "public"."basic_field" VALUES (666, 'detail', 'detail', 2, '', 1, 1614684337);
INSERT INTO "public"."basic_field" VALUES (667, 'geo_sip_country_en', 'geo_sip_country_en', 1, '', 1, 1614684337);
INSERT INTO "public"."basic_field" VALUES (668, 'user', 'user', 2, '', 1, 1614684337);
INSERT INTO "public"."basic_field" VALUES (669, 'host', 'host', 2, '', 1, 1614684337);
INSERT INTO "public"."basic_field" VALUES (670, 'geo_sip_country_cn', 'geo_sip_country_cn', 1, '', 1, 1614684337);
INSERT INTO "public"."basic_field" VALUES (671, 'to', 'to', 2, '', 1, 1614684337);
INSERT INTO "public"."basic_field" VALUES (672, 'method', 'method', 1, '', 1, 1614684337);
INSERT INTO "public"."basic_field" VALUES (673, 'threat_set', 'threat_set', 2, '', 1, 1614684337);
INSERT INTO "public"."basic_field" VALUES (675, 'threat_type', 'threat_type', 1, '', 1, 1614684337);
INSERT INTO "public"."basic_field" VALUES (676, 'plain', 'plain', 2, '', 1, 1614684337);
INSERT INTO "public"."basic_field" VALUES (677, 'sample_direction', 'sample_direction', 1, '', 1, 1614684337);
INSERT INTO "public"."basic_field" VALUES (660, 'geo_dip_country_cn', 'geo_dip_country_cn', 1, '', 1, 1614684337);
INSERT INTO "public"."basic_field" VALUES (678, 'geo_dip_country_en', 'geo_dip_country_en', 1, '', 1, 1614684337);
INSERT INTO "public"."basic_field" VALUES (638, 'log_id', '日志ID', 1, '', 1, 1614684337);
INSERT INTO "public"."basic_field" VALUES (674, 'receive_time', '接收时间', 8, '', 1, 1614684337);
INSERT INTO "public"."basic_field" VALUES (635, 'found_time', '日志请求时间', 8, '', 1, 1614684337);
INSERT INTO "public"."basic_field" VALUES (7, 'data', '日志内容', 2, '', 1, 1614684337);
INSERT INTO "public"."basic_field" VALUES (3, 'event_type', '事件类型', 1, '', 1, 1614684337);
INSERT INTO "public"."basic_field" VALUES (4, 'event_subtype', '事件子类', 1, '', 1, 1614684337);
INSERT INTO "public"."basic_field" VALUES (1, 'source', '来源', 1, '', 1, 1614684337);
INSERT INTO "public"."basic_field" VALUES (5, 'dcd_guid', '装置ID', 1, '', 1, 1614684337);
INSERT INTO "public"."basic_field" VALUES (2, 'service_type', '业务类型', 1, '', 1, 1614684337);
INSERT INTO "public"."basic_field" VALUES (6, 'dev_guid', '资产ID', 1, '', 1, 1614684337);
INSERT INTO "public"."basic_field" VALUES (679, 'device_id', '厂站ID', 5, '', 1, 1614684337);
INSERT INTO "public"."basic_field" VALUES (680, 'netlink_id', '链路id', 5, '', 1, 1614684337);

-- ----------------------------
-- Primary Key structure for table basic_field
-- ----------------------------
ALTER TABLE "public"."basic_field" ADD CONSTRAINT "basic_field_pkey" PRIMARY KEY ("id");
