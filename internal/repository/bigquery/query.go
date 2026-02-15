package bigquery

const ComparisonQuery = `
SELECT
  -- country_code,
  -- campaign,
  sum(
    CASE
      WHEN day >= DATE_SUB(CURRENT_DATE(), INTERVAL 3 day) THEN installs
      ELSE 0
    END
  ) last_3_days_installs,
  sum(
    CASE
      WHEN day >= DATE_SUB(CURRENT_DATE(), INTERVAL 3 day) THEN cost
      ELSE 0
    END
  ) last_3_days_cost,
  sum(
    CASE
      WHEN day >= DATE_SUB(CURRENT_DATE(), INTERVAL 10 day) and day < DATE_SUB(CURRENT_DATE(), INTERVAL 7 day) THEN installs
      ELSE 0
    END
  ) previous_3_days_installs,
  sum(
    CASE
      WHEN day >= DATE_SUB(CURRENT_DATE(), INTERVAL 10 day) and day < DATE_SUB(CURRENT_DATE(), INTERVAL 7 day) THEN cost
      ELSE 0
    END
  ) previous_3_days_cost,
FROM %s.cost_etl
WHERE
  day >= DATE_SUB(CURRENT_DATE(), INTERVAL 10 day)
-- group by 1, 2
order by 1 desc
`

const DatasetIdQuery = `
select
  admin.app_id_to_dataset_id(app_id) as dataset_id
from
  ` + "`admin.apps-metadata`" + `
where
  deleted_at is null
`
