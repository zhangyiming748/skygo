DROP sequence IF EXISTS "public"."table_log_daily_id_seq";

CREATE SEQUENCE "public"."table_log_daily_id_seq"
INCREMENT 1
MINVALUE  1
MAXVALUE 2147483647
START 1
CACHE 1;


-- ----------------------------
-- Table structure for task_log_daily
-- ----------------------------
DROP TABLE IF EXISTS "public"."task_log_daily";
CREATE TABLE "public"."task_log_daily" (
  "id" int4 NOT NULL DEFAULT nextval('table_log_daily_id_seq'::regclass),
  "task_id" int4 NOT NULL,
  "day" int4 NOT NULL DEFAULT 0,
  "total_count" int8 NOT NULL DEFAULT 0,
  "error_count" int8 NOT NULL
)
;
COMMENT ON COLUMN "public"."task_log_daily"."id" IS '自增主键';
COMMENT ON COLUMN "public"."task_log_daily"."task_id" IS '任务主键';
COMMENT ON COLUMN "public"."task_log_daily"."day" IS '天标识';
COMMENT ON COLUMN "public"."task_log_daily"."total_count" IS '处理记录数';
COMMENT ON COLUMN "public"."task_log_daily"."error_count" IS '错误记录数';
COMMENT ON TABLE "public"."task_log_daily" IS '任务日志按天记录';

-- ----------------------------
-- Primary Key structure for table task_log_daily
-- ----------------------------
ALTER TABLE "public"."task_log_daily" ADD CONSTRAINT "table_log_daily_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Foreign Keys structure for table task_log_daily
-- ----------------------------
ALTER TABLE "public"."task_log_daily" ADD CONSTRAINT "fk_table_log_task_1" FOREIGN KEY ("task_id") REFERENCES "public"."task" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;
