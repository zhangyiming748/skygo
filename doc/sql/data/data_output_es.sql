DROP sequence IF EXISTS "public"."data_output_es_id_seq";

CREATE SEQUENCE "public"."data_output_es_id_seq"
INCREMENT 1
MINVALUE  1
MAXVALUE 2147483647
START 1
CACHE 1;

-- ----------------------------
-- Table structure for data_output_es
-- ----------------------------
DROP TABLE IF EXISTS "public"."data_output_es";
CREATE TABLE "public"."data_output_es" (
  "id" int4 NOT NULL DEFAULT nextval('data_output_es_id_seq'::regclass),
  "name" varchar(255) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "config_es_id" int4 NOT NULL,
  "index_name" varchar(255) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "basic_template_id" int4 NOT NULL,
  "create_time" int4 NOT NULL,
  "uuid" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "category" int2 NOT NULL
)
;
COMMENT ON COLUMN "public"."data_output_es"."id" IS '自增ID';
COMMENT ON COLUMN "public"."data_output_es"."name" IS '数据源名称';
COMMENT ON COLUMN "public"."data_output_es"."config_es_id" IS '表cluster_es的记录主键ID';
COMMENT ON COLUMN "public"."data_output_es"."index_name" IS '索引名称';
COMMENT ON COLUMN "public"."data_output_es"."basic_template_id" IS '数据模板';
COMMENT ON COLUMN "public"."data_output_es"."create_time" IS '创建时间';
COMMENT ON COLUMN "public"."data_output_es"."uuid" IS 'uuid';
COMMENT ON COLUMN "public"."data_output_es"."category" IS '分类：1预定义，2自定义';
COMMENT ON TABLE "public"."data_output_es" IS '存储配置es';

-- ----------------------------
-- Records of data_output_es
-- ----------------------------
INSERT INTO "public"."data_output_es" VALUES (1, '协议事件—SSL加密协商', 1, 'iisop_ssl', 41, 1614915021, '05a24605-fcd7-4c18-9816-1dfca8ba3815', 1);
INSERT INTO "public"."data_output_es" VALUES (2, '协议事件—WEB访问', 1, 'iisop_web', 43, 1614915021, 'd88dd9ac-bdd3-4aa4-8dbf-857a48c92395', 1);
INSERT INTO "public"."data_output_es" VALUES (3, '协议事件—域名解析', 1, 'iisop_dns', 33, 1614915021, 'a4b9c614-0cac-40e4-80e1-6d856139b85d', 1);
INSERT INTO "public"."data_output_es" VALUES (4, '威胁数据', 1, 'iisop_threat', 42, 1614915021, 'e91739b2-66ad-41a7-8e38-d0935c3778c6', 1);
INSERT INTO "public"."data_output_es" VALUES (5, '协议事件—工业关键事件', 1, 'iisop_ics_event', 37, 1614915021, '632ab3a7-4453-4fb7-9192-9fb748154f3c', 1);
INSERT INTO "public"."data_output_es" VALUES (6, '协议事件—邮件行为', 1, 'iisop_mail', 39, 1614915021, '197bf7ef-3d52-4f0d-9092-5b52d3b73148', 1);
INSERT INTO "public"."data_output_es" VALUES (7, '流量—会话关闭', 1, 'iisop_flow_end', 36, 1614915021, '2265a212-0bd7-4117-a174-d14af25166b8', 1);
INSERT INTO "public"."data_output_es" VALUES (8, '白名单数据', 1, 'iisop_white', 44, 1614915021, 'a07b18b2-93bd-460f-bf85-6477214fea25', 1);
INSERT INTO "public"."data_output_es" VALUES (9, '中能行为', 1, 'iisop_zhongneng_behavior', 3, 1614915021, '613f46db-166e-4722-8ebc-a964dfdcb540', 1);


-- ----------------------------
-- Primary Key structure for table data_output_es
-- ----------------------------
ALTER TABLE "public"."data_output_es" ADD CONSTRAINT "_copy_2_copy_2_copy_1" PRIMARY KEY ("id");

-- ----------------------------
-- Foreign Keys structure for table data_output_es
-- ----------------------------
ALTER TABLE "public"."data_output_es" ADD CONSTRAINT "fk_data_output_es_basic_template_1" FOREIGN KEY ("basic_template_id") REFERENCES "public"."basic_template" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;
ALTER TABLE "public"."data_output_es" ADD CONSTRAINT "fk_data_output_es_config_elasticsearch_1" FOREIGN KEY ("config_es_id") REFERENCES "public"."config_es" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;
