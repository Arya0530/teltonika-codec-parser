import os
import time
import traceback

from cronjob import run_cronjob


INTERVAL_SECONDS = int(
    os.getenv("CRON_INTERVAL_SECONDS", "300")
)


def main():
    print(
        f"[Scheduler] CronJob aktif setiap "
        f"{INTERVAL_SECONDS} detik"
    )

    while True:
        try:
            run_cronjob()
        except Exception:
            print("[Scheduler] CronJob gagal:")
            traceback.print_exc()

        time.sleep(INTERVAL_SECONDS)


if __name__ == "__main__":
    main()