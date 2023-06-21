DROP sequence IF EXISTS "public"."engine_task_id_seq";

CREATE SEQUENCE "public"."engine_task_id_seq"
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;


-- ----------------------------
-- Table structure for engine_task
-- ----------------------------
DROP TABLE IF EXISTS "public"."engine_task";
CREATE TABLE "public"."engine_task" (
  "id" int4 NOT NULL DEFAULT nextval('engine_task_id_seq'::regclass),
  "engine_task_id" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "start_time" int4 NOT NULL,
  "end_time" int4 NOT NULL,
  "task_id" int4 NOT NULL,
  "final_json" json NOT NULL,
  "total_count" int8 NOT NULL DEFAULT 0,
  "error_count" int8 NOT NULL DEFAULT 0
)
;

-- ----------------------------
-- Indexes structure for table engine_task
-- ----------------------------
CREATE UNIQUE INDEX "unique_engine_task_id" ON "public"."engine_task" USING btree (
  "engine_task_id" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table engine_task
-- ----------------------------
ALTER TABLE "public"."engine_task" ADD CONSTRAINT "engine_task_pkey" PRIMARY KEY ("id");
