from clickhouse import fetch_clickhouse_raw, mark_as_processed
from aggregate import aggregate_data
from postgres import insert_postgres


def run_cronjob():
    rows = fetch_clickhouse_raw(limit=20)

    print("==============================")
    print(f"DATA BELUM DIPROSES DARI CLICKHOUSE: {len(rows)} rows")

    for row in rows:
        processed = aggregate_data(row)

        print("==============================")
        print("HASIL AGGREGATE CRONJOB")
        print(processed)

        insert_postgres(processed)
        mark_as_processed(row)


if __name__ == "__main__":
    run_cronjob()