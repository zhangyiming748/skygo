create view task_test_case_view as
select
tc.id,
tc.task_id,
t.task_uuid,
tc.test_case_id,
k.`name` as `test_case_name`,
k.module_id,
k.auto_test_level,
tc.test_result_status,
tc.action_status,
k.test_tools,
k.task_param,
s.`name` as `scenario_name`
from task_test_case as tc
left join task as t on t.id = tc.task_id
left join knowledge_test_case as k on k.id = tc.test_case_id
left join knowledge_scenario as s on s.id = k.scenario_id