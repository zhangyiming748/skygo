/*
 Navicat Premium Data Transfer

 Source Server         : 乌江60
 Source Server Type    : PostgreSQL
 Source Server Version : 100006
 Source Host           : localhost:5432
 Source Catalog        : skygo_detection
 Source Schema         : public

 Target Server Type    : PostgreSQL
 Target Server Version : 100006
 File Encoding         : 65001

 Date: 05/03/2021 16:10:10
*/


-- ----------------------------
-- Table structure for task_data_output
-- ----------------------------
DROP TABLE IF EXISTS "public"."task_data_output";
CREATE TABLE "public"."task_data_output" (
  "task_id" int4 NOT NULL,
  "data_output_uuid" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "condition" json
)
;
COMMENT ON COLUMN "public"."task_data_output"."task_id" IS '任务id';
COMMENT ON COLUMN "public"."task_data_output"."data_output_uuid" IS '数据存储uuid';
COMMENT ON COLUMN "public"."task_data_output"."condition" IS '描述了数据存储的条件';
COMMENT ON TABLE "public"."task_data_output" IS '任务中数据到指定的数据存储
记录每个任务输出到哪些存储，条件是哪些';

-- ----------------------------
-- Records of task_data_output
-- ----------------------------
INSERT INTO "public"."task_data_output" VALUES (1, '0ac043d6-5253-4b98-93da-c0c80fc14bcd', '{}');
INSERT INTO "public"."task_data_output" VALUES (1, '05a24605-fcd7-4c18-9816-1dfca8ba3815', '{}');
INSERT INTO "public"."task_data_output" VALUES (2, 'd88dd9ac-bdd3-4aa4-8dbf-857a48c92395', '{}');
INSERT INTO "public"."task_data_output" VALUES (2, '77b35784-b038-4752-b3bb-2941269661cc', '{}');
INSERT INTO "public"."task_data_output" VALUES (3, 'a4b9c614-0cac-40e4-80e1-6d856139b85d', '{}');
INSERT INTO "public"."task_data_output" VALUES (3, '50a6806b-44c7-4404-b555-8644cfe7bea7', '{}');
INSERT INTO "public"."task_data_output" VALUES (4, '94d3419f-899b-4528-b96c-41d36650cf4b', '{}');
INSERT INTO "public"."task_data_output" VALUES (4, 'e91739b2-66ad-41a7-8e38-d0935c3778c6', '{}');
INSERT INTO "public"."task_data_output" VALUES (5, '632ab3a7-4453-4fb7-9192-9fb748154f3c', '{}');
INSERT INTO "public"."task_data_output" VALUES (5, '614bd6d0-54a6-4642-90c6-0efe86ae421a', '{}');
INSERT INTO "public"."task_data_output" VALUES (6, 'b1c8fcb0-681f-4d08-8973-3630dc5d1bff', '{}');
INSERT INTO "public"."task_data_output" VALUES (6, '197bf7ef-3d52-4f0d-9092-5b52d3b73148', '{}');
INSERT INTO "public"."task_data_output" VALUES (7, '2265a212-0bd7-4117-a174-d14af25166b8', '{}');
INSERT INTO "public"."task_data_output" VALUES (7, 'e47eefc4-e5b4-4586-bb25-4e207d53ef12', '{}');
INSERT INTO "public"."task_data_output" VALUES (8, 'a07b18b2-93bd-460f-bf85-6477214fea25', '{}');
INSERT INTO "public"."task_data_output" VALUES (8, 'dd0d112b-e680-468c-b572-3467d9f45abe', '{}');
INSERT INTO "public"."task_data_output" VALUES (9, '6e38f00a-28bd-41cb-a827-50f1c86e6c88', '{}');
INSERT INTO "public"."task_data_output" VALUES (10, '613f46db-166e-4722-8ebc-a964dfdcb540', '{}');
INSERT INTO "public"."task_data_output" VALUES (11, '94d3419f-899b-4528-b96c-41d36650cf4b', '{}');

-- ----------------------------
-- Uniques structure for table task_data_output
-- ----------------------------
ALTER TABLE "public"."task_data_output" ADD CONSTRAINT "td" UNIQUE ("task_id", "data_output_uuid");
COMMENT ON CONSTRAINT "td" ON "public"."task_data_output" IS '任务ID和数据输出ID是唯一';

-- ----------------------------
-- Foreign Keys structure for table task_data_output
-- ----------------------------
ALTER TABLE "public"."task_data_output" ADD CONSTRAINT "fk_task_data_output_task_1" FOREIGN KEY ("task_id") REFERENCES "public"."task" ("id") ON DELETE CASCADE ON UPDATE NO ACTION;
