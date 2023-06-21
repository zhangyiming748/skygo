DROP sequence IF EXISTS "public"."basic_log_group_id_seq";

CREATE SEQUENCE "public"."basic_log_group_id_seq"
INCREMENT 1
MINVALUE  1
MAXVALUE 2147483647
START 1
CACHE 1;

SELECT setval('"public"."basic_log_group_id_seq"', 8, true);

-- ----------------------------
-- Table structure for basic_log_group
-- ----------------------------
DROP TABLE IF EXISTS "public"."basic_log_group";
CREATE TABLE "public"."basic_log_group" (
  "id" int4 NOT NULL DEFAULT nextval('basic_log_group_id_seq'::regclass),
  "name" varchar(255) COLLATE "pg_catalog"."default",
  "parent_id" int4,
  "description" varchar(255) COLLATE "pg_catalog"."default",
  "category" int2,
  "create_time" int4
)
;
COMMENT ON COLUMN "public"."basic_log_group"."id" IS '自增主键';
COMMENT ON COLUMN "public"."basic_log_group"."name" IS '字段名称';
COMMENT ON COLUMN "public"."basic_log_group"."parent_id" IS '字段别名';
COMMENT ON COLUMN "public"."basic_log_group"."description" IS '描述';
COMMENT ON COLUMN "public"."basic_log_group"."category" IS '分类：1预定义，2自定义';
COMMENT ON COLUMN "public"."basic_log_group"."create_time" IS '创建时间';
COMMENT ON TABLE "public"."basic_log_group" IS '日志类型';

-- ----------------------------
-- Records of basic_log_group
-- ----------------------------
INSERT INTO "public"."basic_log_group" VALUES (1, '威胁信息', 0, '威胁信息组', 2, 1614684337);
INSERT INTO "public"."basic_log_group" VALUES (6, '工业行为日志', 0, '工业行为日志组', 2, 1614684337);
INSERT INTO "public"."basic_log_group" VALUES (5, '协议行为', 0, '协议行为组', 2, 1614684337);
INSERT INTO "public"."basic_log_group" VALUES (4, '文件上传日志', 0, '文件上传日志组', 2, 1614684337);
INSERT INTO "public"."basic_log_group" VALUES (3, '异常行为日志', 0, '异常行为日志组', 2, 1614684337);
INSERT INTO "public"."basic_log_group" VALUES (2, '流量日志', 0, '流量日志组', 2, 1614684337);
INSERT INTO "public"."basic_log_group" VALUES (7, '南网章和日志', 0, '南网章和日志组', 2, 1614684337);
INSERT INTO "public"."basic_log_group" VALUES (8, '中能融合日志', 0, '中能融合日志', 2, 1614684337);

-- ----------------------------
-- Primary Key structure for table basic_log_group
-- ----------------------------
ALTER TABLE "public"."basic_log_group" ADD CONSTRAINT "_copy_2_copy_3" PRIMARY KEY ("id");
