package sqltool

import (
	"reflect"
	"regexp"
	"testing"
)

func TestGetSyntaxType(t *testing.T) {
	type args struct {
		sql    string
		parser bool
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "create table 带注释",
			args: args{
				sql:    "-- xxx \n CREATE TABLE `admin_entry_menu` (\n  `create_user` varchar(255) NOT NULL DEFAULT '',\n  `modify_user` varchar(255) NOT NULL DEFAULT '',\n  `gmt_create` datetime NOT NULL,\n  `gmt_modify` datetime NOT NULL,\n  `deleted` int(11) NOT NULL DEFAULT '0',\n  `db_remark` varchar(255) NOT NULL DEFAULT '''',\n  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '自增id',\n  `label` varchar(128) NOT NULL DEFAULT '' COMMENT '标签',\n  `icon` varchar(128) NOT NULL DEFAULT '' COMMENT '图标',\n  `url` varchar(128) NOT NULL DEFAULT '' COMMENT '链接',\n  `type` varchar(128) NOT NULL DEFAULT '' COMMENT '类型,default:默认,custom:自定义',\n  `test` int(4) NOT NULL DEFAULT '0',\n  PRIMARY KEY (`id`)\n) ENGINE=InnoDB AUTO_INCREMENT=14 DEFAULT CHARSET=utf8mb4 COMMENT='备注信息';",
				parser: true,
			},
			want: DDL,
		},
		{
			name: "update 带注释",
			args: args{
				sql:    "-- xxxx \n update hbos_scm_statistic.whc_medication_operation_record set is_deleted = 1\nwhere id  in (10128121,10128122,10128123,10128124,10128125,10128126,10128127,10128128,10128130,10128131,10128132,10128133,10128134,10128135,10128136,10128137,10128222,10128223,10128224,10128225,10128226,10128227,10128228,10128229,10128230,10128234,10128235,10128236,10128237,10128238,10128239,10128240,10128241,10128242,10128243,10128244,10128313,10128314,10128315,10128316,10128319,10128320,10128321,10128322,10128323,10128324,10128325,10128326,10128327,10128328,10128329,10128357,10128358,10128359,10128360,10128362,10128363,10128364,10128365,10128401,10128402,10128403,10128404,10128405,10128406,10128407,10128408,10128409,10128410,10128411,10128417,10128419,10128421,10128425,10128426,10128427,10128428,10128429,10128430,10128431,10128432,10128433,10128434,10128435,10128436,10128437,10128440,10128441,10128442,10128443,10128444,10128445,10128446,10128447,10128449,10128452,10128453,10128454,10128455,10128456,10128457,10128458,10128468,10130042,10130043,10130044,10130045,10130056,10130069,10130070,10130071,10130072,10130073,10130074,10130075,10130076,10130077,10130146,10130147,10130148,10130149,10130150,10130151,10130156,10130157,10130159,10130160,10130161,10130162,10130163,10130164,10130165,10130166,10130167,10130168,10130171,10130172,10130177,10130178,10130179,10130180,10130181,10130182,10130183,10130184,10130185,10130186,10130187,10130188,10130190,10130191,10130192,10130193,10130194,10130195,10130197,10130198,10130199,10130200,10130201,10130202,10130203,10130204,10130205,10130206,10130208,10130209,10130210,10130214,10130215,10130216,10130217,10130220,10130221,10130222,10130223,10130225,10130229,10130230,10130231,10130232,10130235,10130236,10130238,10130239,10130240,10130241,10130242,10130243,10130244,10130245,10130247,10130269,10130270,10130271,10130272,10130273,10130398,10130399,10130400,10130718,10130720,10130721,10130722,10130723,10130725,10130727,10130728,10130734,10130735,10130736,10130737,10130738,10130739,10130740,10130741,10130742,10130743,10130744,10130745,10130746,10130747,10130748,10130749,10130750,10130751,10130752,10130753,10130754,10130755,10130756,10130757,10130758,10130759,10130760,10130761,10130762,10130763,10130764,10130765,10130766,10130767,10130768,10130769,10130770,10130771,10130772,10130773,10130774,10130775,10130776,10130779,10130780,10130781,10130782,10130783,10130784,10130785,10130789,10130790,10130791,10130815,10130821,10130822,10130823,10130824,10130825,10130826,10130827,10130828,10130829,10130830,10130831,10130832,10130833,10130834,10130835,10130836,10130837,10130838,10130839,10130840,10130841,10130842,10130843,10130844,10130845,10130846,10130847,10130848,10130849,10130850,10130851,10130852,10130853,10130854,10130855,10130856,10130857,10130866,10130867,10130868,10130869,10130870,10130871,10130872,10130873,10130874,10130875,10130876,10130880,10130881,10130882,10130883,10130884,10130885)",
				parser: true,
			},
			want: DML,
		},
		{
			name: "create view 带注释",
			args: args{
				sql:    "-- xxx \n create or replace VIEW `v_yj_yg_dabian` AS \nselect `si`.`visit_id` AS `BINGRENZYID`,`hr`.`healthcare_record_no` AS `ZHUYUANHAO`,(`si`.`stool_frequency2`) AS `CISHU`,NULL AS `CELIANGSJD`,`si`.`create_time` AS `JILUSJ`,`si`.`create_by` AS `JILURENID`,`si`.`name` AS `JILURENXM`,`si`.`bed_no` AS `CHUANGWEIHAO`,NULL AS `BINGQUID`,NULL AS `BINGQUMC`,`hr`.`dept_id` AS `KESHIID`,`hr`.`dept_name` AS `KESHIMC`,NULL AS `CELIANGSJ`,`hr`.`campus_id` AS `YUANQUID`,`si`.`id` AS `WAIBUID` from (`hbos_emr`.`nis_nurse_sign_info` `si` left join `hbos_diagnosis_treatment`.`dtc_healthcare_record` `hr` on((`hr`.`id` = `si`.`visit_id`))) where (`si`.`stool_frequency2` >=1)\t -- xxx \n",
				parser: true,
			},
			want: DDL,
		},
		{
			name: "alter 带注释",
			args: args{
				sql:    "-- 护理文书主索引表兼容病情护理记录单\nALTER table nis_template_instance\n    add column `master_record_id` bigint(20) default 0 comment '关联主文书记录ID'",
				parser: true,
			},
			want: DDL,
		},
		{
			name: "insert 带注释1",
			args: args{
				sql:    "-- xxx \n INSERT INTO `dict_data` (`main_code`, `sub_code`, `code_no`, `sort_no`, `code_desc`, `dict_type`,\n                         `dict_explain`, `remark`, `parent_code`, `is_deleted`, `create_by`, `update_by`,\n                         `create_time`, `update_time`)\nVALUES ('E0001', 'E0001.007', '007', 7, '文书目录', '系统功能类型', NULL, NULL, NULL, 0, 1, 1, now(), now())",
				parser: true,
			},
			want: DML,
		},
		{
			name: "insert 带注释2",
			args: args{
				sql:    "-- xxx \n insert INTO `dict_data` (`main_code`, `sub_code`, `code_no`, `sort_no`, `code_desc`, `dict_type`,\n                         `dict_explain`, `remark`, `parent_code`, `is_deleted`, `create_by`, `update_by`,\n                         `create_time`, `update_time`)\nVALUES ('E0001', 'E0001.007', '007', 7, '文书目录', '系统功能类型', NULL, NULL, NULL, 0, 1, 1, now(), now())",
				parser: true,
			},
			want: DML,
		},
		{
			name: "delete 1",
			args: args{
				sql:    "-- xxx \n DELETE FROM `hbos_emr`.`dict_system_parameter_template` WHERE `id`=19;",
				parser: true,
			},
			want: DML,
		},
		{
			name: "delete 2",
			args: args{
				sql:    "-- xxx \n delete FROM `hbos_emr`.`dict_system_parameter_template` WHERE `id`=19;",
				parser: true,
			},
			want: DML,
		},
		{
			name: "create 带注释",
			args: args{
				sql:    " -- xxx\n\n\n\n-- insert\ncreate    TABLE `admin_entry_menu` (\n          `create_user` varchar(255) NOT NULL DEFAULT '',\n          `modify_user` varchar(255) NOT NULL DEFAULT '',\n          `gmt_create` datetime NOT NULL,\n          `gmt_modify` datetime NOT NULL,\n          `deleted` int (11) NOT NULL DEFAULT '0',\n          `db_remark` varchar(255) NOT NULL DEFAULT '''',\n          `id` int (11) NOT NULL AUTO_INCREMENT COMMENT '自增id',\n          `label` varchar(128) NOT NULL DEFAULT '' COMMENT '标签',\n          `icon` varchar(128) NOT NULL DEFAULT '' COMMENT '图标',\n          `url` varchar(128) NOT NULL DEFAULT '' COMMENT '链接',\n          `type` varchar(128) NOT NULL DEFAULT '' COMMENT '类型,default:默认,custom:自定义',\n          `test` int (4) NOT NULL DEFAULT '0',\n          PRIMARY KEY (`id`)\n          ) ENGINE = InnoDB AUTO_INCREMENT = 14 DEFAULT CHARSET = utf8mb4 COMMENT = '备注信息';",
				parser: true,
			},
			want: DDL,
		},
		{
			name: "上下文切换 use",
			args: args{
				sql:    "use xxx;",
				parser: true,
			},
			want: Use,
		},
		{
			name: "query 带注释 as中文 1",
			args: args{
				sql:    "-- xx \n select \n count(1) as 检查已报告\n from \n  dtc_doctor_order_executive_plan plan inner join dtc_doctor_order dor on plan.doctor_order_id = dor.id",
				parser: true,
			},
			want: Select,
		},
		{
			name: "query 带注释 as'中文' 2",
			args: args{
				sql:    "-- xxx \n select \n count(1) as '检查已报告'  --xxxx \n \n from \n  dtc_doctor_order_executive_plan plan inner join dtc_doctor_order dor on plan.doctor_order_id = dor.id",
				parser: true,
			},
			want: Select,
		},
		{
			name: "query 带注释 as'中文' 3",
			args: args{
				sql:    "-- xxx \n SELECT    '总计' as '状态',\n          count(1) as '数量'\nfrom      dtc_doctor_order_executive_plan\nwhere     org_id = 20012004\nand       campus_id = 30012005\nand       is_deleted = 0\nand       create_time >= '2022-12-04 19:00:00'\nunion all\nselect    '初始状态' as '状态',\n          count(1) as '数量'\nfrom      dtc_doctor_order_executive_plan\nwhere     org_id = 20012004\nand       campus_id = 30012005\nand       create_time >= '2022-12-04 19:00:00'\nand       is_deleted = 0\nand       search_status in ('FH0180.99', 'FH0181.01')\nunion all\nselect    '执行中' as '状态',\n          count(1) as '数量'\nfrom      dtc_doctor_order_executive_plan\nwhere     org_id = 20012004\nand       campus_id = 30012005\nand       create_time >= '2022-12-04 19:00:00'\nand       is_deleted = 0\nand       search_status in (\n          'FH0180.55',\n          'FH0180.11',\n          'FH0181.02',\n          'FH0180.229',\n          'FH0180.16',\n          'FH0180.01',\n          'FH0180.03',\n          'FH0180.299'\n          )\nunion all\nselect    '已执行' as '状态',\n          count(1) as '数量'\nfrom      dtc_doctor_order_executive_plan\nwhere     org_id = 20012004\nand       campus_id = 30012005\nand       create_time >= '2022-12-04 19:00:00'\nand       is_deleted = 0\nand       search_status in ('FH0180.56', 'FH0181.03')\nunion all\nselect    '已停止' as '状态',\n          count(1) as '数量'\nfrom      dtc_doctor_order_executive_plan\nwhere     org_id = 20012004\nand       campus_id = 30012005\nand       create_time >= '2022-12-04 19:00:00'\nand       is_deleted = 0\nand       search_status in ('FH0181.04')\nunion all\nselect    '已作废' as '状态',\n          count(1) as '数量'\nfrom      dtc_doctor_order_executive_plan\nwhere     org_id = 20012004\nand       campus_id = 30012005\nand       create_time >= '2022-12-04 19:00:00'\nand       is_deleted = 0\nand       search_status in ('FH0180.98')\nunion all\nselect    '其它状态' as '状态',\n          count(1) as '数量'\nfrom      dtc_doctor_order_executive_plan\nwhere     org_id = 20012004\nand       campus_id = 30012005\nand       create_time >= '2022-12-04 19:00:00'\nand       is_deleted = 0\nand       search_status not in (\n          'FH0180.99',\n          'FH0181.01',\n          'FH0180.01',\n          'FH0180.55',\n          'FH0180.11',\n          'FH0181.02',\n          'FH0180.229',\n          'FH0180.56',\n          'FH0181.03',\n          'FH0181.04',\n          'FH0180.98',\n          'FH0180.55',\n          'FH0180.16',\n          'FH0180.03'\n          )",
				parser: true,
			},
			want: Select,
		},
		{
			name: "create database",
			args: args{
				sql:    "create database qqq",
				parser: true,
			},
			want: DbDDL,
		},
		{
			name: "drop table",
			args: args{
				sql:    "drop table qqq",
				parser: true,
			},
			want: DDL,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetSyntaxType(tt.args.sql, tt.args.parser); got != tt.want {
				t.Errorf("GetSyntaxType() = %v, want %v", SyntaxTypeDesc[got], SyntaxTypeDesc[tt.want])
			}
		})
	}
}

func TestZZ(t *testing.T) {
	sensitiveSQLRule := "^create database|^drop database|^rename|^create view|^create or replace view"
	sqlContent := "create or replace view `v_yj_yg_dabian` AS \nselect `si`.`visit_id` AS `BINGRENZYID`,`hr`.`healthcare_record_no` AS `ZHUYUANHAO`,(`si`.`stool_frequency2`) AS `CISHU`,NULL AS `CELIANGSJD`,`si`.`create_time` AS `JILUSJ`,`si`.`create_by` AS `JILURENID`,`si`.`name` AS `JILURENXM`,`si`.`bed_no` AS `CHUANGWEIHAO`,NULL AS `BINGQUID`,NULL AS `BINGQUMC`,`hr`.`dept_id` AS `KESHIID`,`hr`.`dept_name` AS `KESHIMC`,NULL AS `CELIANGSJ`,`hr`.`campus_id` AS `YUANQUID`,`si`.`id` AS `WAIBUID` from (`hbos_emr`.`nis_nurse_sign_info` `si` left join `hbos_diagnosis_treatment`.`dtc_healthcare_record` `hr` on((`hr`.`id` = `si`.`visit_id`))) where (`si`.`stool_frequency2` >=1)\t"

	re, err := regexp.Compile(sensitiveSQLRule)
	if err != nil {
		t.Error(err)
		return
	}

	t.Log(re.MatchString(sqlContent))
}

func TestGenUpdate2Select(t *testing.T) {
	type args struct {
		sqlContent string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name:    "1",
			args:    args{sqlContent: "update `ivc_invoice_make_task` set make_status = 0,system_remark=' where 20230303 重新开票bylly' where id = 10369757 and system_remark = 'where'; "},
			want:    "select count(*) from ivc_invoice_make_task where id = 10369757 and system_remark = 'where'",
			wantErr: false,
		},
		{
			name:    "2",
			args:    args{sqlContent: "update \n  \t `ivc_invoice_make_task` set make_status = 0,system_remark='20230303 重新开票bylly' where id = 10369757; "},
			want:    "select count(*) from ivc_invoice_make_task where id = 10369757",
			wantErr: false,
		},
		{
			name:    "3",
			args:    args{sqlContent: "update  hbos_settlement.stc_bill a\nleft join (\n  select bill_id, count(*) cnt ,sum(case when calc_status in(2,-1) then 1 else 0 end ) calcCnt, sum(origin_amount) as originAmount, sum(billing_amount) billingAmount   from hbos_settlement.stc_bill_detail\n  where 1=1\n  and is_deleted = 0\n  group by  bill_id\n) t on a.id = t.bill_id\nset a.system_remark =  '重症监护室开了ICU治疗单元排斥费用删除-更新账单计算数量 ',\n  a.detail_count = ifnull(t.cnt,0),\n   a.detail_calculated_count = ifnull(t.calcCnt,0),\n   a.bill_amount =  ifnull(t.billingAmount,0),\n   a.update_time = now()\nwhere a.id in (\n10480000\n)\n;"},
			want:    "select count(*) from hbos_settlement.stc_bill as a left join (select bill_id, count(*) as cnt, sum(case when calc_status in (2, -1) then 1 else 0 end) as calcCnt, sum(origin_amount) as originAmount, sum(billing_amount) as billingAmount from hbos_settlement.stc_bill_detail where 1 = 1 and is_deleted = 0 group by bill_id) as t on a.id = t.bill_id where a.id in (10480000)",
			wantErr: false,
		},
		{
			name:    "4",
			args:    args{sqlContent: "update hbos_settlement.stc_settle_trade_step set is_deleted = 1 ,update_time = '2023-05-06 59:59:59', system_remark='历史数据清理，20230506' where org_id = 20034010\t"},
			want:    "select count(*) from hbos_settlement.stc_settle_trade_step where org_id = 20034010",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenUpdate2Select(tt.args.sqlContent)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenUpdate2Select() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GenUpdate2Select() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_updateToSelect(t *testing.T) {
	type args struct {
		sql string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name:    "1",
			args:    args{sql: "update `aaa.ivc_invoice_make_task` set make_status = 0,system_remark=' whe' where id = 10369757; "},
			want:    "select count(*) from `aaa.ivc_invoice_make_task` where id = 10369757",
			wantErr: false,
		},
		{
			name:    "2",
			args:    args{sql: "update `aaa.ivc_invoice_make_task` set make_status = 0,system_remark=' whe' where id = 10369757 and system_remark='where'; "},
			want:    "select count(*) from `aaa.ivc_invoice_make_task` where id = 10369757 and system_remark = 'where'",
			wantErr: false,
		},
		{
			name:    "2",
			args:    args{sql: "delete from table_task where status in (4,5) and processable_time <= '2023-07-30 23:59:59';"},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := genUpdate2SelectV2(tt.args.sql)
			if (err != nil) != tt.wantErr {
				t.Errorf("updateToSelect() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("updateToSelect() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUpdateContainField(t *testing.T) {
	type args struct {
		sql   string
		key   string
		value string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "1",
			args: args{
				sql: "update xxx set id = 1, system_remark='xxx' where 1 = 1",
				key: "system_remark",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "2",
			args: args{
				sql: "update xxx set id = 1, system_remark1='xxx' where 1 = 1",
				key: "system_remark",
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "3",
			args: args{
				sql: "delete from xxx where 1 = 1",
				key: "system_remark",
			},
			want:    true,
			wantErr: false,
		},
		{
			name:    "5",
			args:    args{sql: "select \n if(substring(json_search(medicare.payResult,'all','个人支付'),4,2) REGEXP ']' = 1,JSON_EXTRACT(medicare.payResult,CONCAT('$[',substring(json_search(medicare.payResult,'one','个人支付'),4,1),'].amount')),JSON_EXTRACT(medicare.payResult,CONCAT('$[',substring(json_search(medicare.payResult,'one','个人支付'),4,2),'].amount'))) AS grzf\n-- SUM(if(substring(json_search(medicare.payResult,'all','非医保结算范围费用总额'),4,2) REGEXP ']' = 1,JSON_EXTRACT(medicare.payResult,CONCAT('$[',substring(json_search(medicare.payResult,'one','非医保结算范围费用总额'),4,1),'].amount')),JSON_EXTRACT(medicare.payResult,CONCAT('$[',substring(json_search(medicare.payResult,'one','非医保结算范围费用总额'),4,2),'].amount')))\n--    ) AS '非医保结算范围费用总额',\n-- SUM(case when writeoff_type = 'MEDICARE' then writeoff_amount else 0 END) AS ybzfhj,\n-- SUM(medicare.writeoff_amount) zje\nFROM table111"},
			want:    false,
			wantErr: true,
		},
		{
			name: "6",
			args: args{
				sql:   "update xxx set id = 1, system_remark1=now() where 1 = 1",
				key:   "system_remark1",
				value: "now()",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "7",
			args: args{
				sql:   "update xxx set id = 1, system_remark1 = now() where 1 = 1",
				key:   "system_remark1",
				value: "now()",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "8",
			args: args{
				sql:   "update xxx set id = 1, system_remark1 = now1() where 1 = 1",
				key:   "system_remark1",
				value: "now()",
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CheckUpdateContainField(tt.args.sql, tt.args.key, tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckUpdateContainField() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CheckUpdateContainField() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetQueryFields(t *testing.T) {
	type args struct {
		sql string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name:    "1",
			args:    args{sql: "select xxx1, xxxx as x111, *, id as id1 from user0001"},
			want:    []string{"xxx1", "x111", "*", "id1"},
			wantErr: false,
		},
		{
			name:    "2",
			args:    args{sql: "select \n\n      t1.id ,t1.sku_id,t2.before_retail_price, t2.before_purchase_price\n\n      from whc_check_form_detail t1\n\n      inner join \n      (\n        select    sku_id,\n                  before_retail_price,\n                  before_purchase_price\n        from      whc_price_adjust_form_detail\n        where     price_adjust_no = '2023022100022'\n      ) t2 on t1.sku_id = t2.sku_id\n      where t1.check_no = '2023030700014'"},
			want:    []string{"id", "sku_id", "before_retail_price", "before_purchase_price"},
			wantErr: false,
		},
		{
			name:    "3",
			args:    args{sql: "select xxx1., xxxx as x111, *, id as id1 from user0001"},
			want:    []string{},
			wantErr: false,
		},
		{
			name:    "3",
			args:    args{sql: "explain select xxx1., xxxx as x111, *, id as id1 from user0001"},
			want:    []string{},
			wantErr: false,
		},
		{
			name:    "4",
			args:    args{sql: "select\ndm.id as IIID,\n       case(dm.medication_type)\n           when 'FH0101.10.01' then '西药'\n           when 'FH0101.10.02' then '中成药'\n           when 'FH0101.10.03' then '草药' else '未知' end as '药品类型',\n       dm.name as '药品名称',\n       dm.specification as '规格',\n       dm.manufacturer_name as '生产厂商',\n       dm.purchase_price as '采购价',\n       dm.retail_price as '零售价',\ndm.identification_code as '医保国家编码',\n       concat('1',dm.large_package_unit,'=',dm.small_large_package_ratio,dm.small_package_unit) as '包装',\n       IF( dm.`status` = 0, '是', '否' ) AS '是否停用',\n       dm.small_package_unit as '小包装',\n       dm.large_package_unit as '大包装',\n       (select sum(total_num) from hbos_inventory.itc_inventory where is_deleted=0 and org_id=20012004 and campus_id=30012005 and sku_id=dm.id) \n       as '库存数',\n            dm.license_number as '批准文号'\n\n from hbos_medication.mc_dict_medication\ndm  where dm.is_deleted=0\nand dm.is_deleted=0 and dm.org_id=20012004\n  order by dm.id desc;"},
			want:    []string{"IIID", "药品类型", "药品名称", "规格", "生产厂商", "采购价", "零售价", "医保国家编码", "包装", "是否停用", "小包装", "大包装", "库存数", "批准文号"},
			wantErr: false,
		},
		{
			name:    "4",
			args:    args{sql: "select aaa as '阿达是的' from table1"},
			want:    []string{"阿达是的"},
			wantErr: false,
		},
		{
			name:    "",
			args:    args{sql: "select id, name as '姓名', q as '群', age as '年龄' from user0001 where id < 10"},
			want:    []string{"id", "姓名", "群", "年龄"},
			wantErr: false,
		},
		{
			name:    "",
			args:    args{sql: "select id, name qqq, q  '群', age '年龄' from user0001 where id < 10"},
			want:    []string{"id", "qqq", "群", "年龄"},
			wantErr: false,
		},
		{
			name:    "",
			args:    args{sql: "select id, name as q1, q as q2, age as q3 from user0001 where id < 10"},
			want:    []string{"id", "q1", "q2", "q3"},
			wantErr: false,
		},
		{
			name:    "",
			args:    args{sql: "select mdm.name                                                         as   '药品名称',\n\nmmrb.purchase_platform_code    as    '采购平台id',\n\n             (select dict_value_name\n\n              from apex_dic.dic_common_value dcv\n\n              where dcv.org_id = mdm.org_id\n\n                and dict_value = mdm.medication_type)                  as            '药品类型',\n\n             dose_form_name                                                   as   '剂型',\n\n             specification                                                    as   '规格',\n\n             manufacturer_name                                                as   '生产厂商',\n\n             CONCAT('1', large_package_unit, '=', small_large_package_ratio,\n\n                    small_package_unit)                                       as     '包装',\n\n             CONCAT('1', small_package_unit, '=', JSON_EXTRACT(dose_info, '$[1][0].dose[1]'),\n\n                    REPLACE(JSON_EXTRACT(dose_info, '$[1][0].doseUnit'), '', '')) as '剂量1',\n\n             CONCAT('1', small_package_unit, '=', JSON_EXTRACT(dose_info, '$[1][1].dose[1]'),\n\n                    REPLACE(JSON_EXTRACT(dose_info, '$[1][1].doseUnit'), '', ''))  as '剂量2',\n\n             CONCAT('1', small_package_unit, '=', JSON_EXTRACT(dose_info, '$[1][2].dose[1]'),\n\n                    REPLACE(JSON_EXTRACT(dose_info, '$[1][2].doseUnit'), '', '')) as '剂量3',\n\n             CONCAT(JSON_EXTRACT(mmrb.order_limit, '$.adultDoseLimit.onceDoseNum[1]'),\n\n                    REPLACE(JSON_EXTRACT(mmrb.order_limit, '$.adultDoseLimit.onceDoseNumUnit'), '',\n\n                            ''))                                              as   '成人默认单次剂量',\n\n             CONCAT(JSON_EXTRACT(mmrb.order_limit, '$.adultDoseLimit.onceDoseNumMax[1]'),\n\n                    REPLACE(JSON_EXTRACT(mmrb.order_limit, '$.adultDoseLimit.onceDoseNumMaxUnit'), '',\n\n                            ''))                                              as   '成人一次量上限',\n\n             CONCAT(JSON_EXTRACT(mmrb.order_limit, '$.adultDoseLimit.dailyDoseNumMax[1]'),\n\n                    REPLACE(JSON_EXTRACT(mmrb.order_limit, '$.adultDoseLimit.dailyDoseNumMaxUnit'), '',\n\n                            ''))                                              as   '成人一日量上限',\n\n             REPLACE(JSON_EXTRACT(mdmc.default_usage, '$.name'), '', '')     as   '默认用法',\n\n             REPLACE(JSON_EXTRACT(mdmc.default_frequency, '$.name'), '', '') as   '默认频次',\n\n             mdmc.special_drugs_name                                          as     '毒理分类',\n\n             mdmc.high_risk_level_name                                        as     '高危等级',\n\n             mdmc.antimicrobial_grade_name                                    as     '抗菌药等级',\n\n             if(mdmc.is_skin_test is null, mdmc.is_skin_test,\n\n                IF(mdmc.is_skin_test = 1, '是', '否'))                           as    '是否皮试',\n\n             mdmc.antitumor_name                                              as     '抗肿瘤药物',\n\n             CONCAT(mdm.retail_price, '元/', large_package_unit)              as      '零售价',\n\n             CONCAT(mdm.purchase_price, '元/', large_package_unit)            as      '采购价',\n\n             if(mdm.is_basic_medication is null, mdm.is_basic_medication,\n\n                if(mdm.is_basic_medication = 1, '是', '否'))                    as   '是否基药',\n\n             if(mmrb.is_cent_purchase is null, mmrb.is_cent_purchase,\n\n                if(mmrb.is_cent_purchase = 1, '是', '否'))                      as   '是否集采药品',\n\n             if(mmrb.is_refrigerate is null, mmrb.is_refrigerate,\n\n                if(mmrb.is_refrigerate = 1, '是', '否'))                        as   '是否冷藏药品',\n\n             REPLACE(JSON_EXTRACT(mmrb.round_info, '$.outpatientRound.name'), '',\n\n                     '')                                                      as     '门、急诊取整方式',\n\n             REPLACE(JSON_EXTRACT(mmrb.round_info, '$.inpatientRound.name'), '',\n\n                     '')                                                      as     '住院取整方式',\n\n             REPLACE(JSON_EXTRACT(mmrb.round_info, '$.observationRound.name'), '',\n\n                     '')                                                     as      '留观取整方式',\n\n             IF(REPLACE(JSON_EXTRACT(mdm.features, '$.definedDailyDose'), '',\n\n                        '') is null, REPLACE(JSON_EXTRACT(mdm.features, '$.definedDailyDose'), '',\n\n                                             ''), REPLACE(JSON_EXTRACT(mdm.features, '$.definedDailyDose[1]'), '',\n\n                                                          ''))            as         'DDD值',\n\n             mdm.identification_code                                    as           '医保国家编码',\n\n             if(mdm.status = 1, '否', '是')                                     as   '是否停用',\n\n             mdm.id                                                           as   id\n\n      from hbos_medication_1.mc_dict_medication_common mdmc\n\n               inner join hbos_medication_1.mc_dict_medication mdm on mdmc.id = mdm.medication_common_id\n\n               inner join hbos_medication_1.mc_medication_rule_biz mmrb on mdm.id = mmrb.medication_id\n\n      where  medication_type in ('FH0101.10.01', 'FH0101.10.02')\n\n    and mdm.is_deleted = 0\n\n    and mdm.status = 1;\n\n\n\n"},
			want:    []string{"药品名称", "采购平台id", "药品类型", "剂型", "规格", "生产厂商", "包装", "剂量1", "剂量2", "剂量3", "成人默认单次剂量", "成人一次量上限", "成人一日量上限", "默认用法", "默认频次", "毒理分类", "高危等级", "抗菌药等级", "是否皮试", "抗肿瘤药物", "零售价", "采购价", "是否基药", "是否集采药品", "是否冷藏药品", "门、急诊取整方式", "住院取整方式", "留观取整方式", "DDD值", "医保国家编码", "是否停用", "id"},
			wantErr: false,
		},
		{
			name:    "",
			args:    args{sql: "SELECT  c.*\nFROM    stc_bill b\nINNER JOIN stc_settle_bill c\nON      c.id = b.settle_bill_id\nWHERE   b.patient_name = '钟铁军'\nAND     b.settle_status = 'SETTLED'\nAND     SUBSTR(b.settle_time,1,10) = '2023-05-29'\nAND     b.is_deleted = 0\nAND     c.is_deleted = 0\nAND     b.org_id = '20012004'\n;"},
			want:    []string{"*"},
			wantErr: false,
		},
		{
			name:    "",
			args:    args{sql: `select id id中文 from app_app`},
			want:    []string{"id中文"},
			wantErr: false,
		},
		{
			name:    "",
			args:    args{sql: `explain select id id中文 from app_app`},
			want:    []string{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetQueryFields(tt.args.sql)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetQueryFields() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetQueryFields() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEscapeSql(t *testing.T) {
	type args struct {
		sql string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "1",
			args: args{sql: "urn:hl7-org:v3 ..\\..\\cdaschemas\\sdschemas\\CDA'.xsd"},
			want: "urn:hl7-org:v3 ..\\\\..\\\\cdaschemas\\\\sdschemas\\\\CDA\\'.xsd",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EscapeSql(tt.args.sql); got != tt.want {
				t.Errorf("EscapeSql() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetTablesFromSelectStatement(t *testing.T) {
	type args struct {
		query string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name:    "1",
			args:    args{query: "SELECT t1.*, t2.name FROM db1.table1 AS t1 LEFT JOIN db2.table2 AS t2 ON t1.id = t2.id"},
			want:    []string{"db1.table1", "db2.table2"},
			wantErr: false,
		},
		{
			name:    "1",
			args:    args{query: "SELECT *, t2.name FROM db1.table1 LEFT JOIN db2.table2 AS t2 ON id = t2.id"},
			want:    []string{"db1.table1", "db2.table2"},
			wantErr: false,
		},
		{
			name:    "",
			args:    args{query: "SELECT  c.*\nFROM    stc_bill b\nINNER JOIN stc_settle_bill c\nON      c.id = b.settle_bill_id\nWHERE   b.patient_name = '钟铁军'\nAND     b.settle_status = 'SETTLED'\nAND     SUBSTR(b.settle_time,1,10) = '2023-05-29'\nAND     b.is_deleted = 0\nAND     c.is_deleted = 0\nAND     b.org_id = '20012004'\n;"},
			want:    []string{"stc_bill", "stc_settle_bill"},
			wantErr: false,
		},
		{
			name:    "",
			args:    args{query: "SELECT  c.*\nFROM    stc_bill b\nINNER JOIN stc_settle_bill c\nON      c.id = b.settle_bill_id\nWHERE   b.patient_name = '钟铁军'\nAND     b.settle_status = 'SETTLED'\nAND     SUBSTR(b.settle_time,1,10) = '2023-05-29'\nAND     b.is_deleted = 0\nAND     c.is_deleted = 0\nAND     b.org_id = '20012004'\n;"},
			want:    []string{"stc_bill", "stc_settle_bill"},
			wantErr: false,
		},
		{
			name:    "",
			args:    args{query: "select * from a"},
			want:    []string{"a"},
			wantErr: false,
		},
		{
			name:    "",
			args:    args{query: "select    i.external_code 外部编码,\n\n          i.id 医嘱项目ID,\n\n          i.management_category_code as 管理类目,\n\n          mc.name 管理类目名称,\n\n          i.name 项目名称,\n\n          i.code 医嘱项目,\n\n          i.unit as 单位,\n\n          i.name as 医嘱缩写,\n\n          (\n\n          select    group_concat(p.country_price_code)\n\n          from      hbos_basedata.hsc_service_item_relation r\n\n          left join hbos_basedata.hsc_charge_item c on c.is_deleted = 0\n\n          and       c.org_id = 20034010\n\n          and       c.code = r.charge_item_code\n\n          and       c.name not like '%2部位%'\n\n          and       c.name not like '%3部位%'\n\n          left join hbos_basedata.hsc_standard_price p on p.code = c.price_code\n\n          and       p.is_deleted = 0\n\n          and       p.org_id = 20034010\n\n          where     r.is_deleted = 0\n\n          and       r.org_id = 20034010\n\n          and       r.main_body_type = 1\n\n          and       r.main_body_code = i.code\n\n          ) 物价编码,\n\n          (\n\n          select    group_concat(c.name)\n\n          from      hbos_basedata.hsc_service_item_relation r\n\n          left join hbos_basedata.hsc_charge_item c on c.is_deleted = 0\n\n          and       c.org_id = 20034010\n\n          and       c.code = r.charge_item_code\n\n          and       c.name not like '%2部位%'\n\n          and       c.name not like '%3部位%'\n\n          left join hbos_basedata.hsc_standard_price p on p.code = c.price_code\n\n          and       p.is_deleted = 0\n\n          and       p.org_id = 20034010\n\n          where     r.is_deleted = 0\n\n          and       r.org_id = 20034010\n\n          and       r.main_body_type = 1\n\n          and       r.main_body_code = i.code\n\n          ) 收费项目名,\n\n          (\n\n          select    group_concat(p.price)\n\n          from      hbos_basedata.hsc_service_item_relation r\n\n          left join hbos_basedata.hsc_charge_item c on c.is_deleted = 0\n\n          and       c.org_id = 20034010\n\n          and       c.code = r.charge_item_code\n\n          and       c.name not like '%2部位%'\n\n          and       c.name not like '%3部位%'\n\n          left join hbos_basedata.hsc_standard_price p on p.code = c.price_code\n\n          and       p.is_deleted = 0\n\n          and       p.org_id = 20034010\n\n          where     r.is_deleted = 0\n\n          and       r.org_id = 20034010\n\n          and       r.main_body_type = 1\n\n          and       r.main_body_code = i.code\n\n          ) 收费项目名单价,\n\n          (\n\n          select    group_concat(r.quantity)\n\n          from      hbos_basedata.hsc_service_item_relation r\n\n          left join hbos_basedata.hsc_charge_item c on c.is_deleted = 0\n\n          and       c.org_id = 20034010\n\n          and       c.code = r.charge_item_code\n\n          and       c.name not like '%2部位%'\n\n          and       c.name not like '%3部位%'\n\n          left join hbos_basedata.hsc_standard_price p on p.code = c.price_code\n\n          and       p.is_deleted = 0\n\n          and       p.org_id = 20034010\n\n          where     r.is_deleted = 0\n\n          and       r.org_id = 20034010\n\n          and       r.main_body_type = 1\n\n          and       r.main_body_code = i.code\n\n          ) 收费项目数量,\n\n          (\n\n          select    sum(p.price * r.quantity)\n\n          from      hbos_basedata.hsc_service_item_relation r\n\n          left join hbos_basedata.hsc_charge_item c on c.is_deleted = 0\n\n          and       c.org_id = 20034010\n\n          and       c.code = r.charge_item_code\n\n          and       c.name not like '%2部位%'\n\n          and       c.name not like '%3部位%'\n\n          left join hbos_basedata.hsc_standard_price p on p.code = c.price_code\n\n          and       p.is_deleted = 0\n\n          and       p.org_id = 20034010\n\n          where     r.is_deleted = 0\n\n          and       r.org_id = 20034010\n\n          and       r.main_body_type = 1\n\n          and       r.main_body_code = i.code\n\n          ) 收费项目总金额,\n\n          (\n\n          select    group_concat(ai.code) 附加服务项目code\n\n          from      hbos_basedata.hsc_service_item_addition a\n\n          left join hbos_basedata.hsc_service_item ai on ai.is_deleted = 0\n\n          and       ai.code = a.service_code\n\n          and       ai.org_id = 20034010\n\n          where     a.is_deleted = 0\n\n          and       a.service_type = 'FH0134.02'\n\n          and       a.main_body_code = i.code\n\n          and       a.main_body_type = '1'\n\n          and       a.org_id = 20034010\n\n          ) 附加服务项目code,\n\n          (\n\n          select    group_concat(ai.name) 附加服务项目名称\n\n          from      hbos_basedata.hsc_service_item_addition a\n\n          left join hbos_basedata.hsc_service_item ai on ai.is_deleted = 0\n\n          and       ai.code = a.service_code\n\n          and       ai.org_id = 20034010\n\n          where     a.is_deleted = 0\n\n          and       a.service_type = 'FH0134.02'\n\n          and       a.main_body_code = i.code\n\n          and       a.main_body_type = '1'\n\n          and       a.org_id = 20034010\n\n          ) 附加服务项目名称,\n\n          (\n\n          select    group_concat(a.quantity) 附加服务项目数量\n\n          from      hbos_basedata.hsc_service_item_addition a\n\n          left join hbos_basedata.hsc_service_item ai on ai.is_deleted = 0\n\n          and       ai.code = a.service_code\n\n          and       ai.org_id = 20034010\n\n          where     a.is_deleted = 0\n\n          and       a.service_type = 'FH0134.02'\n\n          and       a.main_body_code = i.code\n\n          and       a.main_body_type = '1'\n\n          and       a.org_id = 20034010\n\n          ) 附加服务项目数量,\n\n          (\n\n          select    group_concat(pp.price) 附加服务项目价格\n\n          from      hbos_basedata.hsc_service_item_addition a\n\n          left join hbos_basedata.hsc_service_item ai on ai.is_deleted = 0\n\n          and       ai.code = a.service_code\n\n          left join hbos_basedata.hsc_service_item_relation rr on rr.is_deleted = 0\n\n          and       rr.org_id = 20034010\n\n          and       rr.main_body_type = 1\n\n          and       rr.main_body_code = ai.code\n\n          left join hbos_basedata.hsc_charge_item cc on cc.is_deleted = 0\n\n          and       cc.org_id = 20034010\n\n          and       cc.code = rr.charge_item_code\n\n          left join hbos_basedata.hsc_standard_price pp on pp.code = cc.price_code\n\n          and       pp.is_deleted = 0\n\n          and       pp.org_id = 20034010\n\n          and       ai.org_id = 20034010\n\n          where     a.is_deleted = 0\n\n          and       a.service_type = 'FH0134.02'\n\n          and       a.main_body_code = i.code\n\n          and       a.main_body_type = '1'\n\n          and       a.org_id = 20034010\n\n          ) 附加服务_项目单价,\n\n          (\n\n          select    group_concat(rr.quantity) 附加服务项目_收费数量\n\n          from      hbos_basedata.hsc_service_item_addition a\n\n          left join hbos_basedata.hsc_service_item ai on ai.is_deleted = 0\n\n          and       ai.code = a.service_code\n\n          left join hbos_basedata.hsc_service_item_relation rr on rr.is_deleted = 0\n\n          and       rr.org_id = 20034010\n\n          and       rr.main_body_type = 1\n\n          and       rr.main_body_code = ai.code\n\n          left join hbos_basedata.hsc_charge_item cc on cc.is_deleted = 0\n\n          and       cc.org_id = 20034010\n\n          and       cc.code = rr.charge_item_code\n\n          left join hbos_basedata.hsc_standard_price pp on pp.code = cc.price_code\n\n          and       pp.is_deleted = 0\n\n          and       pp.org_id = 20034010\n\n          and       ai.org_id = 20034010\n\n          where     a.is_deleted = 0\n\n          and       a.service_type = 'FH0134.02'\n\n          and       a.main_body_code = i.code\n\n          and       a.main_body_type = '1'\n\n          and       a.org_id = 20034010\n\n          ) 附加服务项目_收费数量,\n\n          (\n\n          select    sum(rr.quantity * pp.price) 附加服务项目_收费数量\n\n          from      hbos_basedata.hsc_service_item_addition a\n\n          left join hbos_basedata.hsc_service_item ai on ai.is_deleted = 0\n\n          and       ai.code = a.service_code\n\n          left join hbos_basedata.hsc_service_item_relation rr on rr.is_deleted = 0\n\n          and       rr.org_id = 20034010\n\n          and       rr.main_body_type = 1\n\n          and       rr.main_body_code = ai.code\n\n          left join hbos_basedata.hsc_charge_item cc on cc.is_deleted = 0\n\n          and       cc.org_id = 20034010\n\n          and       cc.code = rr.charge_item_code\n\n          left join hbos_basedata.hsc_standard_price pp on pp.code = cc.price_code\n\n          and       pp.is_deleted = 0\n\n          and       pp.org_id = 20034010\n\n          and       ai.org_id = 20034010\n\n          where     a.is_deleted = 0\n\n          and       a.service_type = 'FH0134.02'\n\n          and       a.main_body_code = i.code\n\n          and       a.main_body_type = '1'\n\n          and       a.org_id = 20034010\n\n          ) 附加服务项目_总金额,\n\n          (\n\n          select    group_concat(s.id)\n\n          from      hbos_basedata.hsc_service_item_addition a1\n\n          left join hbos_basedata.mtc_sku s on s.is_deleted = 0\n\n          and       s.org_id = 20034010\n\n          and       s.code = a1.service_code\n\n          where     a1.is_deleted = 0\n\n          and       a1.service_type = 'FH0134.03'\n\n          and       a1.main_body_code = i.code\n\n          and       a1.main_body_type = '1'\n\n          and       a1.org_id = 20034010\n\n          ) 附加材料id,\n\n          (\n\n          select    group_concat(s.name)\n\n          from      hbos_basedata.hsc_service_item_addition a1\n\n          left join hbos_basedata.mtc_sku s on s.is_deleted = 0\n\n          and       s.org_id = 20034010\n\n          and       s.code = a1.service_code\n\n          where     a1.is_deleted = 0\n\n          and       a1.service_type = 'FH0134.03'\n\n          and       a1.main_body_code = i.code\n\n          and       a1.main_body_type = '1'\n\n          and       a1.org_id = 20034010\n\n          ) 附加材料名称,\n\n          (\n\n          select    group_concat(a1.quantity)\n\n          from      hbos_basedata.hsc_service_item_addition a1\n\n          left join hbos_basedata.mtc_sku s on s.is_deleted = 0\n\n          and       s.org_id = 20034010\n\n          and       s.code = a1.service_code\n\n          where     a1.is_deleted = 0\n\n          and       a1.service_type = 'FH0134.03'\n\n          and       a1.main_body_code = i.code\n\n          and       a1.main_body_type = '1'\n\n          and       a1.org_id = 20034010\n\n          ) 附加材料数量,\n\n          (\n\n          select    group_concat(s.min_count_sell_price)\n\n          from      hbos_basedata.hsc_service_item_addition a1\n\n          left join hbos_basedata.mtc_sku s on s.is_deleted = 0\n\n          and       s.org_id = 20034010\n\n          and       s.code = a1.service_code\n\n          where     a1.is_deleted = 0\n\n          and       a1.service_type = 'FH0134.03'\n\n          and       a1.main_body_code = i.code\n\n          and       a1.main_body_type = '1'\n\n          and       a1.org_id = 20034010\n\n          ) 附加材料价格,\n\n          (\n\n          select    sum(s.min_count_sell_price * a1.quantity)\n\n          from      hbos_basedata.hsc_service_item_addition a1\n\n          left join hbos_basedata.mtc_sku s on s.is_deleted = 0\n\n          and       s.org_id = 20034010\n\n          and       s.code = a1.service_code\n\n          where     a1.is_deleted = 0\n\n          and       a1.service_type = 'FH0134.03'\n\n          and       a1.main_body_code = i.code\n\n          and       a1.main_body_type = '1'\n\n          and       a1.org_id = 20034010\n\n          ) 附加材料金额,\n\n          (\n\n          select    group_concat(dd.name)\n\n          from      hbos_basedata.hsc_service_item_addition_relation rrr\n\n          left join hbos_basedata.hsc_service_item_addition_group a on a.is_deleted = 0\n\n          and       a.id = rrr.relation_addition_id\n\n          and       a.org_id = 20034010\n\n          left join hbos_basedata.hsc_service_item_addition aaa on aaa.is_deleted = 0\n\n          and       aaa.org_id = 20034010\n\n          and       aaa.main_body_code = a.group_code\n\n          and       aaa.service_type = 'FH0134.01'\n\n          and       aaa.main_body_type = 3\n\n          left join hbos_medication.mc_dict_medication dd on dd.is_deleted = 0 -- and dd.org_id =  20034010\n\n          and       aaa.service_type = 'FH0134.01'\n\n          and       dd.id = aaa.service_code\n\n          where     rrr.is_deleted = 0\n\n          and       rrr.main_body_code = i.code\n\n          and       rrr.main_body_type = '1'\n\n          and       rrr.org_id = 20034010\n\n          ) 带药,\n\n          (\n\n          select    sum(dd.retail_price * aaa.quantity)\n\n          from      hbos_basedata.hsc_service_item_addition_relation rrr\n\n          left join hbos_basedata.hsc_service_item_addition_group a on a.is_deleted = 0\n\n          and       a.id = rrr.relation_addition_id\n\n          and       a.org_id = 20034010\n\n          left join hbos_basedata.hsc_service_item_addition aaa on aaa.is_deleted = 0\n\n          and       aaa.org_id = 20034010\n\n          and       aaa.main_body_code = a.group_code\n\n          and       aaa.service_type = 'FH0134.01'\n\n          and       aaa.main_body_type = 3\n\n          left join hbos_medication.mc_dict_medication dd on dd.is_deleted = 0 -- and dd.org_id =  20034010\n\n          and       aaa.service_type = 'FH0134.01'\n\n          and       dd.id = aaa.service_code\n\n          where     rrr.is_deleted = 0\n\n          and       rrr.main_body_code = i.code\n\n          and       rrr.main_body_type = '1'\n\n          and       rrr.org_id = 20034010\n\n          ) 带药金额\n\nfrom      hbos_basedata.hsc_service_item i\n\ninner     join hbos_basedata.hsc_management_category mc on mc.is_deleted = 0\n\nand       mc.code = i.management_category_code\n\nand       mc.org_id = 20034010\n\nwhere     i.is_deleted = 0\n\nand       i.org_id = 20034010 --  AND i.name='无痛胃镜（含色素内镜）'\n\nAND       management_category_code LIKE '%FH0101.01%'\n\norder by  mc.code"},
			want:    []string{"hbos_basedata.hsc_service_item", "hbos_basedata.hsc_management_category"},
			wantErr: false,
		},
		{
			name:    "",
			args:    args{query: "select id from (  select a.* from app_app as a limit 1  ) as ab"},
			want:    []string{"app_app"},
			wantErr: false,
		},
		{
			name: "",
			args: args{query: `
select    项目id,
          项目code,
          项目名称,
          参考价格,
          sum(数量 * 单价) as 价格
from      (
          select    a.id as 项目id,
                    a.code as 项目code,
                    a.name as 项目名称,
                    a.reference_price as 参考价格,
                    '物价' as '类型',
                    c.code as 收费code,
                    c.name as 收费项目名称,
                    b.quantity as 数量,
                    d.price as 单价
          from      hbos_basedata_11.hsc_service_item a
          join      hbos_basedata_11.hsc_service_item_relation b on a.code = b.main_body_code
          and       b.is_deleted = 0
          join      hbos_basedata_11.hsc_charge_item c on b.charge_item_code = c.code
          and       c.is_deleted = 0
          join      hbos_basedata_11.hsc_standard_price d on d.code = c.price_code
          and       d.is_deleted = 0
          where     a.is_deleted = 0
          and       a.status = 1
          ) as tt
where     1 = 1
group by  项目id,
          项目code,
          项目名称,
          参考价格;
`},
			want: []string{
				"hbos_basedata_11.hsc_service_item",
				"hbos_basedata_11.hsc_service_item_relation",
				"hbos_basedata_11.hsc_charge_item",
				"hbos_basedata_11.hsc_standard_price",
			},
			wantErr: false,
		},
		{
			name: "",
			args: args{query: `
select b.id  通用名id,c.identification_code  国家医保编码 from hbos_medication.mc_dict_medication_classify a 
left join hbos_medication.mc_dict_medication_common b on b.medication_classify_id = a.id and b.is_deleted = 0
left join hbos_medication.mc_dict_medication c on c.medication_common_id = b.id and c.is_deleted = 0
where  a.is_deleted = 0 
`},
			want: []string{
				"hbos_medication.mc_dict_medication_classify",
				"hbos_medication.mc_dict_medication_common",
				"hbos_medication.mc_dict_medication",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetTablesFromSelectStatement(tt.args.query)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTablesFromSelectStatement() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetTablesFromSelectStatement() got = %v, want %v", got, tt.want)
				for _, s := range got {
					t.Errorf(s)
				}
			}
		})
	}
}

func TestCheckInsertContainField(t *testing.T) {
	type args struct {
		sql string
		key string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "",
			args: args{
				sql: "INSERT INTO mytable (id, name) VALUES (1, 'John')",
				key: "name",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "",
			args: args{
				sql: "INSERT INTO mytable (id, name) VALUES (1, 'John'), (2, '')",
				key: "name",
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "",
			args: args{
				sql: "INSERT Ignore INTO hbos_diagnosis_treatment.dtc_diagnosis_dict (id, name, type, sub_type, pinyin, icd_code,\n                                                                is_contagious,\n                                                                create_time, update_time, code, description, memo,\n                                                                search_code,\n                                                                morphology_code, enable_state, for_outpatient,\n                                                                for_hospitalization, for_emergency,\n                                                                for_emergency_observation,\n                                                                for_emergency_rescue, is_deleted, wbm, category, org_id,\n                                                                product_scope, system_remark)\nselect id + 1000000,\n       name,\n       type,\n       sub_type,\n       pinyin,\n       icd_code,\n       is_contagious,\n       '2023-05-18 00:00:00',\n       '2023-05-18 00:00:00',\n       code,\n       description,\n       memo,\n       search_code,\n       morphology_code,\n       enable_state,\n       for_outpatient,\n       for_hospitalization,\n       for_emergency,\n       for_emergency_observation,\n       for_emergency_rescue,\n       is_deleted,\n       wbm,\n       category,\n       20058003,\n       product_scope,\n       system_remark\nfrom dtc_diagnosis_dict\nwhere org_id = 20001003\n  and id = 10002001;",
				key: "system_remark",
			},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CheckInsertContainField(tt.args.sql, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckInsertContainField() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CheckInsertContainField() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parserType(t *testing.T) {
	type args struct {
		sql string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "",
			args: args{sql: "use xxx;"},
			want: Use,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parserType(tt.args.sql); got != tt.want {
				t.Errorf("parserType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_warpSqlChinese(t *testing.T) {
	type args struct {
		sql string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{
			name: "",
			args: args{sql: "select id, name qqq, q  群, age '年龄' from user0001 where id < 10"},
			want: "select id, name qqq, q  '群', age '年龄' from user0001 where id < 10",
		},
		{
			name: "",
			args: args{sql: `
select mdm.name                                                         as   '药品名称',

mmrb.purchase_platform_code    as    '采购平台id',

             (select dict_value_name

              from apex_dic.dic_common_value dcv

              where dcv.org_id = mdm.org_id

                and dict_value = mdm.medication_type)                  as            '药品类型',

             dose_form_name                                                   as   '剂型',

             specification                                                    as   '规格',

             manufacturer_name                                                as   '生产厂商',

             CONCAT('1', large_package_unit, '=', small_large_package_ratio,

                    small_package_unit)                                       as     '包装',

             CONCAT('1', small_package_unit, '=', JSON_EXTRACT(dose_info, '$[1][0].dose[1]'),

                    REPLACE(JSON_EXTRACT(dose_info, '$[1][0].doseUnit'), '', '')) as '剂量1',

             CONCAT('1', small_package_unit, '=', JSON_EXTRACT(dose_info, '$[1][1].dose[1]'),

                    REPLACE(JSON_EXTRACT(dose_info, '$[1][1].doseUnit'), '', ''))  as '剂量2',

             CONCAT('1', small_package_unit, '=', JSON_EXTRACT(dose_info, '$[1][2].dose[1]'),

                    REPLACE(JSON_EXTRACT(dose_info, '$[1][2].doseUnit'), '', '')) as '剂量3',

             CONCAT(JSON_EXTRACT(mmrb.order_limit, '$.adultDoseLimit.onceDoseNum[1]'),

                    REPLACE(JSON_EXTRACT(mmrb.order_limit, '$.adultDoseLimit.onceDoseNumUnit'), '',

                            ''))                                              as   '成人默认单次剂量',

             CONCAT(JSON_EXTRACT(mmrb.order_limit, '$.adultDoseLimit.onceDoseNumMax[1]'),

                    REPLACE(JSON_EXTRACT(mmrb.order_limit, '$.adultDoseLimit.onceDoseNumMaxUnit'), '',

                            ''))                                              as   '成人一次量上限',

             CONCAT(JSON_EXTRACT(mmrb.order_limit, '$.adultDoseLimit.dailyDoseNumMax[1]'),

                    REPLACE(JSON_EXTRACT(mmrb.order_limit, '$.adultDoseLimit.dailyDoseNumMaxUnit'), '',

                            ''))                                              as   '成人一日量上限',

             REPLACE(JSON_EXTRACT(mdmc.default_usage, '$.name'), '', '')     as   '默认用法',

             REPLACE(JSON_EXTRACT(mdmc.default_frequency, '$.name'), '', '') as   '默认频次',

             mdmc.special_drugs_name                                          as     '毒理分类',

             mdmc.high_risk_level_name                                        as     '高危等级',

             mdmc.antimicrobial_grade_name                                    as     '抗菌药等级',

             if(mdmc.is_skin_test is null, mdmc.is_skin_test,

                IF(mdmc.is_skin_test = 1, '是', '否'))                           as    '是否皮试',

             mdmc.antitumor_name                                              as     '抗肿瘤药物',

             CONCAT(mdm.retail_price, '元/', large_package_unit)              as      '零售价',

             CONCAT(mdm.purchase_price, '元/', large_package_unit)            as      '采购价',

             if(mdm.is_basic_medication is null, mdm.is_basic_medication,

                if(mdm.is_basic_medication = 1, '是', '否'))                    as   '是否基药',

             if(mmrb.is_cent_purchase is null, mmrb.is_cent_purchase,

                if(mmrb.is_cent_purchase = 1, '是', '否'))                      as   '是否集采药品',

             if(mmrb.is_refrigerate is null, mmrb.is_refrigerate,

                if(mmrb.is_refrigerate = 1, '是', '否'))                        as   '是否冷藏药品',

             REPLACE(JSON_EXTRACT(mmrb.round_info, '$.outpatientRound.name'), '',

                     '')                                                      as     '门、急诊取整方式',

             REPLACE(JSON_EXTRACT(mmrb.round_info, '$.inpatientRound.name'), '',

                     '')                                                      as     '住院取整方式',

             REPLACE(JSON_EXTRACT(mmrb.round_info, '$.observationRound.name'), '',

                     '')                                                     as      '留观取整方式',

             IF(REPLACE(JSON_EXTRACT(mdm.features, '$.definedDailyDose'), '',

                        '') is null, REPLACE(JSON_EXTRACT(mdm.features, '$.definedDailyDose'), '',

                                             ''), REPLACE(JSON_EXTRACT(mdm.features, '$.definedDailyDose[1]'), '',

                                                          ''))            as         'DDD值',

             mdm.identification_code                                    as           '医保国家编码',

             if(mdm.status = 1, '否', '是')                                     as   '是否停用',

             mdm.id                                                           as   id

      from hbos_medication_1.mc_dict_medication_common mdmc

               inner join hbos_medication_1.mc_dict_medication mdm on mdmc.id = mdm.medication_common_id

               inner join hbos_medication_1.mc_medication_rule_biz mmrb on mdm.id = mmrb.medication_id

      where  medication_type in ('FH0101.10.01', 'FH0101.10.02')

    and mdm.is_deleted = 0

    and mdm.status = 1;



`},
			want: "",
		},
		{
			name: "",
			args: args{sql: "SELECT  c.*\nFROM    stc_bill b\nINNER JOIN stc_settle_bill c\nON      c.id = b.settle_bill_id\nWHERE   b.patient_name = '钟铁军'\nAND     b.settle_status = 'SETTLED'\nAND     SUBSTR(b.settle_time,1,10) = '2023-05-29'\nAND     b.is_deleted = 0\nAND     c.is_deleted = 0\nAND     b.org_id = '20012004'\n;"},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := warpSqlChinese(tt.args.sql)
			t.Log(got)
		})
	}
}

func TestParseTableName(t *testing.T) {
	type args struct {
		table string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		want1   string
		wantErr bool
	}{
		{
			name:    "",
			args:    args{table: "matrix.app_app"},
			want:    "matrix",
			want1:   "app_app",
			wantErr: false,
		},
		{
			name:    "",
			args:    args{table: "matrix.app_app as a"},
			want:    "matrix",
			want1:   "app_app",
			wantErr: false,
		},
		{
			name:    "",
			args:    args{table: "app_app"},
			want:    "",
			want1:   "app_app",
			wantErr: false,
		},
		{
			name:    "",
			args:    args{table: "app_app as a"},
			want:    "",
			want1:   "app_app",
			wantErr: false,
		},
		{
			name:    "",
			args:    args{table: "(select a.* from app_app as a limit 1)"},
			want:    "",
			want1:   "",
			wantErr: true,
		},
		{
			name:    "",
			args:    args{table: "MATRIX as a"},
			want:    "",
			want1:   "MATRIX",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := ParseTableName(tt.args.table)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseTableName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseTableName() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("ParseTableName() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_warpSqlChinese2(t *testing.T) {
	type args struct {
		sql string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "",
			args: args{sql: "select id, app_code 应用Code1, app_name 应用Code名称2 from app_app limit 1"},
			want: " select id, app_code '应用Code1', app_name '应用Code名称2' from app_app limit 1",
		},
		{
			name: "",
			args: args{sql: "select    i.external_code 外部编码,\n\n          i.id 医嘱项目ID,\n\n          i.management_category_code as 管理类目,\n\n          mc.name 管理类目名称,\n\n          i.name 项目名称,\n\n          i.code 医嘱项目,\n\n          i.unit as 单位,\n\n          i.name as 医嘱缩写,\n\n          (\n\n          select    group_concat(p.country_price_code)\n\n          from      hbos_basedata.hsc_service_item_relation r\n\n          left join hbos_basedata.hsc_charge_item c on c.is_deleted = 0\n\n          and       c.org_id = 20034010\n\n          and       c.code = r.charge_item_code\n\n          and       c.name not like '%2部位%'\n\n          and       c.name not like '%3部位%'\n\n          left join hbos_basedata.hsc_standard_price p on p.code = c.price_code\n\n          and       p.is_deleted = 0\n\n          and       p.org_id = 20034010\n\n          where     r.is_deleted = 0\n\n          and       r.org_id = 20034010\n\n          and       r.main_body_type = 1\n\n          and       r.main_body_code = i.code\n\n          ) 物价编码,\n\n          (\n\n          select    group_concat(c.name)\n\n          from      hbos_basedata.hsc_service_item_relation r\n\n          left join hbos_basedata.hsc_charge_item c on c.is_deleted = 0\n\n          and       c.org_id = 20034010\n\n          and       c.code = r.charge_item_code\n\n          and       c.name not like '%2部位%'\n\n          and       c.name not like '%3部位%'\n\n          left join hbos_basedata.hsc_standard_price p on p.code = c.price_code\n\n          and       p.is_deleted = 0\n\n          and       p.org_id = 20034010\n\n          where     r.is_deleted = 0\n\n          and       r.org_id = 20034010\n\n          and       r.main_body_type = 1\n\n          and       r.main_body_code = i.code\n\n          ) 收费项目名,\n\n          (\n\n          select    group_concat(p.price)\n\n          from      hbos_basedata.hsc_service_item_relation r\n\n          left join hbos_basedata.hsc_charge_item c on c.is_deleted = 0\n\n          and       c.org_id = 20034010\n\n          and       c.code = r.charge_item_code\n\n          and       c.name not like '%2部位%'\n\n          and       c.name not like '%3部位%'\n\n          left join hbos_basedata.hsc_standard_price p on p.code = c.price_code\n\n          and       p.is_deleted = 0\n\n          and       p.org_id = 20034010\n\n          where     r.is_deleted = 0\n\n          and       r.org_id = 20034010\n\n          and       r.main_body_type = 1\n\n          and       r.main_body_code = i.code\n\n          ) 收费项目名单价,\n\n          (\n\n          select    group_concat(r.quantity)\n\n          from      hbos_basedata.hsc_service_item_relation r\n\n          left join hbos_basedata.hsc_charge_item c on c.is_deleted = 0\n\n          and       c.org_id = 20034010\n\n          and       c.code = r.charge_item_code\n\n          and       c.name not like '%2部位%'\n\n          and       c.name not like '%3部位%'\n\n          left join hbos_basedata.hsc_standard_price p on p.code = c.price_code\n\n          and       p.is_deleted = 0\n\n          and       p.org_id = 20034010\n\n          where     r.is_deleted = 0\n\n          and       r.org_id = 20034010\n\n          and       r.main_body_type = 1\n\n          and       r.main_body_code = i.code\n\n          ) 收费项目数量,\n\n          (\n\n          select    sum(p.price * r.quantity)\n\n          from      hbos_basedata.hsc_service_item_relation r\n\n          left join hbos_basedata.hsc_charge_item c on c.is_deleted = 0\n\n          and       c.org_id = 20034010\n\n          and       c.code = r.charge_item_code\n\n          and       c.name not like '%2部位%'\n\n          and       c.name not like '%3部位%'\n\n          left join hbos_basedata.hsc_standard_price p on p.code = c.price_code\n\n          and       p.is_deleted = 0\n\n          and       p.org_id = 20034010\n\n          where     r.is_deleted = 0\n\n          and       r.org_id = 20034010\n\n          and       r.main_body_type = 1\n\n          and       r.main_body_code = i.code\n\n          ) 收费项目总金额,\n\n          (\n\n          select    group_concat(ai.code) 附加服务项目code\n\n          from      hbos_basedata.hsc_service_item_addition a\n\n          left join hbos_basedata.hsc_service_item ai on ai.is_deleted = 0\n\n          and       ai.code = a.service_code\n\n          and       ai.org_id = 20034010\n\n          where     a.is_deleted = 0\n\n          and       a.service_type = 'FH0134.02'\n\n          and       a.main_body_code = i.code\n\n          and       a.main_body_type = '1'\n\n          and       a.org_id = 20034010\n\n          ) 附加服务项目code,\n\n          (\n\n          select    group_concat(ai.name) 附加服务项目名称\n\n          from      hbos_basedata.hsc_service_item_addition a\n\n          left join hbos_basedata.hsc_service_item ai on ai.is_deleted = 0\n\n          and       ai.code = a.service_code\n\n          and       ai.org_id = 20034010\n\n          where     a.is_deleted = 0\n\n          and       a.service_type = 'FH0134.02'\n\n          and       a.main_body_code = i.code\n\n          and       a.main_body_type = '1'\n\n          and       a.org_id = 20034010\n\n          ) 附加服务项目名称,\n\n          (\n\n          select    group_concat(a.quantity) 附加服务项目数量\n\n          from      hbos_basedata.hsc_service_item_addition a\n\n          left join hbos_basedata.hsc_service_item ai on ai.is_deleted = 0\n\n          and       ai.code = a.service_code\n\n          and       ai.org_id = 20034010\n\n          where     a.is_deleted = 0\n\n          and       a.service_type = 'FH0134.02'\n\n          and       a.main_body_code = i.code\n\n          and       a.main_body_type = '1'\n\n          and       a.org_id = 20034010\n\n          ) 附加服务项目数量,\n\n          (\n\n          select    group_concat(pp.price) 附加服务项目价格\n\n          from      hbos_basedata.hsc_service_item_addition a\n\n          left join hbos_basedata.hsc_service_item ai on ai.is_deleted = 0\n\n          and       ai.code = a.service_code\n\n          left join hbos_basedata.hsc_service_item_relation rr on rr.is_deleted = 0\n\n          and       rr.org_id = 20034010\n\n          and       rr.main_body_type = 1\n\n          and       rr.main_body_code = ai.code\n\n          left join hbos_basedata.hsc_charge_item cc on cc.is_deleted = 0\n\n          and       cc.org_id = 20034010\n\n          and       cc.code = rr.charge_item_code\n\n          left join hbos_basedata.hsc_standard_price pp on pp.code = cc.price_code\n\n          and       pp.is_deleted = 0\n\n          and       pp.org_id = 20034010\n\n          and       ai.org_id = 20034010\n\n          where     a.is_deleted = 0\n\n          and       a.service_type = 'FH0134.02'\n\n          and       a.main_body_code = i.code\n\n          and       a.main_body_type = '1'\n\n          and       a.org_id = 20034010\n\n          ) 附加服务_项目单价,\n\n          (\n\n          select    group_concat(rr.quantity) 附加服务项目_收费数量\n\n          from      hbos_basedata.hsc_service_item_addition a\n\n          left join hbos_basedata.hsc_service_item ai on ai.is_deleted = 0\n\n          and       ai.code = a.service_code\n\n          left join hbos_basedata.hsc_service_item_relation rr on rr.is_deleted = 0\n\n          and       rr.org_id = 20034010\n\n          and       rr.main_body_type = 1\n\n          and       rr.main_body_code = ai.code\n\n          left join hbos_basedata.hsc_charge_item cc on cc.is_deleted = 0\n\n          and       cc.org_id = 20034010\n\n          and       cc.code = rr.charge_item_code\n\n          left join hbos_basedata.hsc_standard_price pp on pp.code = cc.price_code\n\n          and       pp.is_deleted = 0\n\n          and       pp.org_id = 20034010\n\n          and       ai.org_id = 20034010\n\n          where     a.is_deleted = 0\n\n          and       a.service_type = 'FH0134.02'\n\n          and       a.main_body_code = i.code\n\n          and       a.main_body_type = '1'\n\n          and       a.org_id = 20034010\n\n          ) 附加服务项目_收费数量,\n\n          (\n\n          select    sum(rr.quantity * pp.price) 附加服务项目_收费数量\n\n          from      hbos_basedata.hsc_service_item_addition a\n\n          left join hbos_basedata.hsc_service_item ai on ai.is_deleted = 0\n\n          and       ai.code = a.service_code\n\n          left join hbos_basedata.hsc_service_item_relation rr on rr.is_deleted = 0\n\n          and       rr.org_id = 20034010\n\n          and       rr.main_body_type = 1\n\n          and       rr.main_body_code = ai.code\n\n          left join hbos_basedata.hsc_charge_item cc on cc.is_deleted = 0\n\n          and       cc.org_id = 20034010\n\n          and       cc.code = rr.charge_item_code\n\n          left join hbos_basedata.hsc_standard_price pp on pp.code = cc.price_code\n\n          and       pp.is_deleted = 0\n\n          and       pp.org_id = 20034010\n\n          and       ai.org_id = 20034010\n\n          where     a.is_deleted = 0\n\n          and       a.service_type = 'FH0134.02'\n\n          and       a.main_body_code = i.code\n\n          and       a.main_body_type = '1'\n\n          and       a.org_id = 20034010\n\n          ) 附加服务项目_总金额,\n\n          (\n\n          select    group_concat(s.id)\n\n          from      hbos_basedata.hsc_service_item_addition a1\n\n          left join hbos_basedata.mtc_sku s on s.is_deleted = 0\n\n          and       s.org_id = 20034010\n\n          and       s.code = a1.service_code\n\n          where     a1.is_deleted = 0\n\n          and       a1.service_type = 'FH0134.03'\n\n          and       a1.main_body_code = i.code\n\n          and       a1.main_body_type = '1'\n\n          and       a1.org_id = 20034010\n\n          ) 附加材料id,\n\n          (\n\n          select    group_concat(s.name)\n\n          from      hbos_basedata.hsc_service_item_addition a1\n\n          left join hbos_basedata.mtc_sku s on s.is_deleted = 0\n\n          and       s.org_id = 20034010\n\n          and       s.code = a1.service_code\n\n          where     a1.is_deleted = 0\n\n          and       a1.service_type = 'FH0134.03'\n\n          and       a1.main_body_code = i.code\n\n          and       a1.main_body_type = '1'\n\n          and       a1.org_id = 20034010\n\n          ) 附加材料名称,\n\n          (\n\n          select    group_concat(a1.quantity)\n\n          from      hbos_basedata.hsc_service_item_addition a1\n\n          left join hbos_basedata.mtc_sku s on s.is_deleted = 0\n\n          and       s.org_id = 20034010\n\n          and       s.code = a1.service_code\n\n          where     a1.is_deleted = 0\n\n          and       a1.service_type = 'FH0134.03'\n\n          and       a1.main_body_code = i.code\n\n          and       a1.main_body_type = '1'\n\n          and       a1.org_id = 20034010\n\n          ) 附加材料数量,\n\n          (\n\n          select    group_concat(s.min_count_sell_price)\n\n          from      hbos_basedata.hsc_service_item_addition a1\n\n          left join hbos_basedata.mtc_sku s on s.is_deleted = 0\n\n          and       s.org_id = 20034010\n\n          and       s.code = a1.service_code\n\n          where     a1.is_deleted = 0\n\n          and       a1.service_type = 'FH0134.03'\n\n          and       a1.main_body_code = i.code\n\n          and       a1.main_body_type = '1'\n\n          and       a1.org_id = 20034010\n\n          ) 附加材料价格,\n\n          (\n\n          select    sum(s.min_count_sell_price * a1.quantity)\n\n          from      hbos_basedata.hsc_service_item_addition a1\n\n          left join hbos_basedata.mtc_sku s on s.is_deleted = 0\n\n          and       s.org_id = 20034010\n\n          and       s.code = a1.service_code\n\n          where     a1.is_deleted = 0\n\n          and       a1.service_type = 'FH0134.03'\n\n          and       a1.main_body_code = i.code\n\n          and       a1.main_body_type = '1'\n\n          and       a1.org_id = 20034010\n\n          ) 附加材料金额,\n\n          (\n\n          select    group_concat(dd.name)\n\n          from      hbos_basedata.hsc_service_item_addition_relation rrr\n\n          left join hbos_basedata.hsc_service_item_addition_group a on a.is_deleted = 0\n\n          and       a.id = rrr.relation_addition_id\n\n          and       a.org_id = 20034010\n\n          left join hbos_basedata.hsc_service_item_addition aaa on aaa.is_deleted = 0\n\n          and       aaa.org_id = 20034010\n\n          and       aaa.main_body_code = a.group_code\n\n          and       aaa.service_type = 'FH0134.01'\n\n          and       aaa.main_body_type = 3\n\n          left join hbos_medication.mc_dict_medication dd on dd.is_deleted = 0 -- and dd.org_id =  20034010\n\n          and       aaa.service_type = 'FH0134.01'\n\n          and       dd.id = aaa.service_code\n\n          where     rrr.is_deleted = 0\n\n          and       rrr.main_body_code = i.code\n\n          and       rrr.main_body_type = '1'\n\n          and       rrr.org_id = 20034010\n\n          ) 带药,\n\n          (\n\n          select    sum(dd.retail_price * aaa.quantity)\n\n          from      hbos_basedata.hsc_service_item_addition_relation rrr\n\n          left join hbos_basedata.hsc_service_item_addition_group a on a.is_deleted = 0\n\n          and       a.id = rrr.relation_addition_id\n\n          and       a.org_id = 20034010\n\n          left join hbos_basedata.hsc_service_item_addition aaa on aaa.is_deleted = 0\n\n          and       aaa.org_id = 20034010\n\n          and       aaa.main_body_code = a.group_code\n\n          and       aaa.service_type = 'FH0134.01'\n\n          and       aaa.main_body_type = 3\n\n          left join hbos_medication.mc_dict_medication dd on dd.is_deleted = 0 -- and dd.org_id =  20034010\n\n          and       aaa.service_type = 'FH0134.01'\n\n          and       dd.id = aaa.service_code\n\n          where     rrr.is_deleted = 0\n\n          and       rrr.main_body_code = i.code\n\n          and       rrr.main_body_type = '1'\n\n          and       rrr.org_id = 20034010\n\n          ) 带药金额\n\nfrom      hbos_basedata.hsc_service_item i\n\ninner     join hbos_basedata.hsc_management_category mc on mc.is_deleted = 0\n\nand       mc.code = i.management_category_code\n\nand       mc.org_id = 20034010\n\nwhere     i.is_deleted = 0\n\nand       i.org_id = 20034010 --  AND i.name='无痛胃镜（含色素内镜）'\n\nAND       management_category_code LIKE '%FH0101.01%'\n\norder by  mc.code"},
			want: "",
		},
		{
			name: "",
			args: args{sql: `
select b.id 通用名id,c.identification_code 国家医保编码 from hbos_medication.mc_dict_medication_classify a 
left join hbos_medication.mc_dict_medication_common b on b.medication_classify_id = a.id and b.is_deleted = 0
left join hbos_medication.mc_dict_medication c on c.medication_common_id = b.id and c.is_deleted = 0
where  a.is_deleted = 0 
`},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := warpSqlChinese2(tt.args.sql)
			t.Log(got)
		})
	}
}

func TestCheckSqlAlterAndGetTableCode(t *testing.T) {
	type args struct {
		sql string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "",
			args: args{
				sql: "alter table qqq.`tableq` add COLUMN ep_flag tinyint(1) DEFAULT 1 COMMENT '表单修改权限 0:关闭 1:开启';",
			},
			want: "tableq",
		},
		{
			name: "",
			args: args{
				sql: "alter table `qqq`.`tableq` add COLUMN ep_flag tinyint(1) DEFAULT 1 COMMENT '表单修改权限 0:关闭 1:开启';",
			},
			want: "tableq",
		},
		{
			name: "",
			args: args{
				sql: "alter table `tableq` add COLUMN ep_flag tinyint(1) DEFAULT 1 COMMENT '表单修改权限 0:关闭 1:开启';",
			},
			want: "tableq",
		},
		{
			name: "",
			args: args{
				sql: "alter table dic_register add COLUMN ep_flag tinyint(1) DEFAULT 1 COMMENT '表单修改权限 0:关闭 1:开启';",
			},
			want: "dic_register",
		},
		{
			name: "",
			args: args{
				sql: "alter table dic_register add COLUMN ep_flag tinyint(1) DEFAULT 1 COMMENT '表单修改权限 0:关闭 1:开启',ALGORITHM=INPLACE, LOCK=NONE;",
			},
			want: "dic_register",
		},
		{
			name: "",
			args: args{
				sql: `
				/** 请在晚上空闲时执行 **/

				alter table dtc_doctor_order
					modify extension mediumtext null comment '扩展信息' ,ALGORITHM = INPLACE,LOCK = NONE;
				`,
			},
			want: "dtc_doctor_order",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CheckSqlAlterAndGetTableCode(tt.args.sql)
			if got != tt.want {
				t.Errorf("CheckSqlAlterAndGetTableCode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetSelectTableInfosFromSelectStatement(t *testing.T) {
	type args struct {
		query string
	}
	tests := []struct {
		name    string
		args    args
		want    []SelectTableInfo
		wantErr bool
	}{
		{
			name: "",
			args: args{
				query: "select * from hbos_medication.mc_dict_medication",
			},
			want: []SelectTableInfo{{
				DbCode:    "hbos_medication",
				TableCode: "mc_dict_medication",
			}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetSelectTableInfosFromSelectStatement(tt.args.query)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetSelectTableInfosFromSelectStatement() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetSelectTableInfosFromSelectStatement() = %v, want %v", got, tt.want)
			}
		})
	}
}
