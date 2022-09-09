\c rakutan

CREATE MATERIALIZED VIEW mat_view_rakutan AS
    SELECT
        *
         ,passed_total::float / register_total::float as rate
    FROM (
             SELECT
                 *
                  ,(SELECT SUM(r) FROM UNNEST(register) r) as register_total -- 過去n年分の履修人数
                  ,(SELECT SUM(p) FROM UNNEST(passed) p) as passed_total -- 過去n年分の単位取得人数
                  ,(SELECT COUNT(c) FROM UNNEST(array_remove(register, NULL)) c) as count -- 過去n年のうち開講回数
             FROM rakutan
         ) r
    WHERE
                register_total / count > 20	-- 平均20人以上受講してる講義
      AND r.faculty_name = '国際高等教育院'
      AND passed_total::float / register_total::float > 0.75 -- 単位取得率が7.5割以上
      AND id < 20000 -- 昨年開講されているもの
    ;

CREATE MATERIALIZED VIEW mat_view_onitan AS
    SELECT
        *
         ,passed_total::float / register_total::float as rate
    FROM (
             SELECT
                 *
                  ,(SELECT SUM(r) FROM UNNEST(register) r) as register_total -- 過去n年分の履修人数
                  ,(SELECT SUM(p) FROM UNNEST(passed) p) as passed_total -- 過去n年分の単位取得人数
                  ,(SELECT COUNT(c) FROM UNNEST(array_remove(register, NULL)) c) as count -- 過去n年のうち開講回数
             FROM rakutan
         ) r
    WHERE
                register_total / count > 2	-- 平均3人以上受講してる講義
      AND passed_total::float / register_total::float < 0.35 -- 単位取得率が3.5割以下
      AND id < 20000 -- 昨年開講されているもの
    ;