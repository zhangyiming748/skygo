DROP sequence IF EXISTS "public"."config_field_map_id_seq" ;

CREATE SEQUENCE "public"."config_field_map_id_seq"
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;


-- ----------------------------
-- Table structure for config_field_map
-- ----------------------------
DROP TABLE IF EXISTS "public"."config_field_map";
CREATE TABLE "public"."config_field_map" (
  "id" int4 NOT NULL DEFAULT nextval('config_field_map_id_seq'::regclass),
  "name" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "json_config" varchar(1500) COLLATE "pg_catalog"."default"
)
;
COMMENT ON COLUMN "public"."config_field_map"."id" IS '主键id';
COMMENT ON COLUMN "public"."config_field_map"."name" IS '数据字典名称';
COMMENT ON COLUMN "public"."config_field_map"."json_config" IS '配置，用一个json存储';
COMMENT ON TABLE "public"."config_field_map" IS '字段配置';

-- ----------------------------
-- Records of config_field_map
-- ----------------------------

-- ----------------------------
-- Primary Key structure for table config_field_map
-- ----------------------------
ALTER TABLE "public"."config_field_map" ADD CONSTRAINT "_copy_5" PRIMARY KEY ("id");
