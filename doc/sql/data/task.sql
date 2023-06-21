DROP sequence IF EXISTS "public"."task_id_seq";

CREATE SEQUENCE "public"."task_id_seq"
INCREMENT 1
MINVALUE  1
MAXVALUE 2147483647
START 1
CACHE 1;


-- ----------------------------
-- Table structure for task
-- ----------------------------
DROP TABLE IF EXISTS "public"."task";
CREATE TABLE "public"."task" (
  "id" int4 NOT NULL DEFAULT nextval('task_id_seq'::regclass),
  "name" varchar(255) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "description" varchar(255) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "data_input_uuid" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "rule_id" int4 NOT NULL,
  "create_time" int4 NOT NULL,
  "total_count" int8 NOT NULL DEFAULT 0,
  "error_count" int8 NOT NULL DEFAULT 0,
  "final_json" json,
  "status" int2 NOT NULL DEFAULT 1,
  "run_engine_task_id" varchar(255) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "run_start_time" int4 NOT NULL DEFAULT 0,
  "run_end_time" int4 NOT NULL DEFAULT 0
)
;
COMMENT ON COLUMN "public"."task"."id" IS '自增主键';
COMMENT ON COLUMN "public"."task"."name" IS '任务名称';
COMMENT ON COLUMN "public"."task"."description" IS '描述';
COMMENT ON COLUMN "public"."task"."data_input_uuid" IS '数据采集id';
COMMENT ON COLUMN "public"."task"."rule_id" IS '解析规则id';
COMMENT ON COLUMN "public"."task"."create_time" IS '创建时间';
COMMENT ON COLUMN "public"."task"."total_count" IS '总数';
COMMENT ON COLUMN "public"."task"."error_count" IS '错误数量';
COMMENT ON COLUMN "public"."task"."final_json" IS '给引擎用的json指令';
COMMENT ON COLUMN "public"."task"."status" IS '状态，1草稿，2未运行 3运行中 4异常';
COMMENT ON TABLE "public"."task" IS '任务表';

-- ----------------------------
-- Records of task
-- ----------------------------
INSERT INTO "public"."task" VALUES (1, '协议-SSL加密协商上报', '', 'a925e5c2-f6a3-4a67-a821-6bd2329049c0', 2, 1614684337, 0, 0, '{}', 1, '', 0, 0);
INSERT INTO "public"."task" VALUES (2, '协议-WEB访问上报', '', 'a8408ad5-96bf-4fa3-a6d7-8e257cbba001', 3, 1614684337, 0, 0, '{}', 1, '', 0, 0);
INSERT INTO "public"."task" VALUES (3, '协议-域名解析上报', '', 'd82174a9-5e0d-43a8-a145-78b3596ea3cf', 1, 1614684337, 0, 0, NULL, 1, '', 0, 0);
INSERT INTO "public"."task" VALUES (5, '协议-工业关键事件上报', '', '288b76d5-fbc6-4aa9-99ba-ad8e8eed200b', 7, 1614684337, 0, 0, NULL, 1, '', 0, 0);
INSERT INTO "public"."task" VALUES (4, '协议-威胁数据上报', '', '563de29a-ebde-4cc0-a1ab-f5af80d89993', 8, 1614684337, 0, 0, NULL, 1, '', 0, 0);
INSERT INTO "public"."task" VALUES (6, '协议-邮件行为上报', '', '856cb227-de32-4509-8a55-c0bbd95b65d5', 6, 1614684337, 0, 0, NULL, 1, '', 0, 0);
INSERT INTO "public"."task" VALUES (7, '流量-会话关闭上报', '', '3b79ec48-d83c-4249-b264-e6ebb5b6c8b1', 5, 1614684337, 0, 0, NULL, 1, '', 0, 0);
INSERT INTO "public"."task" VALUES (8, '白名单数据上报', '', 'fcf97c35-d229-446a-841f-5561dee59aab', 4, 1614684337, 0, 0, NULL, 1, '', 0, 0);
INSERT INTO "public"."task" VALUES (9, '资产上报', '', 'f854cdff-2abc-419d-a9ea-50deac72e1b3', 9, 1614684337, 0, 0, NULL, 1, '', 0, 0);
INSERT INTO "public"."task" VALUES (10, '中能融合-行为日志上报', '', 'cc9f6216-a852-48c1-8336-0cc2b63e44a1', 10, 1614684337, 0, 0, NULL, 1, '', 0, 0);
INSERT INTO "public"."task" VALUES (11, '中能融合-警报日志上报', '', '7e6274db-bce4-47c6-ad30-42459f25a87f', 10, 1614684337, 0, 0, NULL, 1, '', 0, 0);

-- ----------------------------
-- Primary Key structure for table task
-- ----------------------------
ALTER TABLE "public"."task" ADD CONSTRAINT "_copy_7" PRIMARY KEY ("id");

-- ----------------------------
-- Foreign Keys structure for table task
-- ----------------------------
ALTER TABLE "public"."task" ADD CONSTRAINT "fk_task_rule_1" FOREIGN KEY ("rule_id") REFERENCES "public"."rule" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;
