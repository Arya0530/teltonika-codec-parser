-- CreateTable
CREATE TABLE "devices" (
    "id" SERIAL NOT NULL,
    "imei" VARCHAR(30) NOT NULL,
    "device_name" VARCHAR(100),
    "is_active" BOOLEAN NOT NULL DEFAULT true,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "devices_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "processed_avl" (
    "id" SERIAL NOT NULL,
    "device_id" INTEGER NOT NULL,
    "timestamp" TIMESTAMP(3) NOT NULL,
    "latitude" DOUBLE PRECISION NOT NULL,
    "longitude" DOUBLE PRECISION NOT NULL,
    "speed" INTEGER NOT NULL,
    "status" VARCHAR(20) NOT NULL,
    "priority" INTEGER NOT NULL,
    "altitude" INTEGER NOT NULL,
    "angle" INTEGER NOT NULL,
    "satellites" INTEGER NOT NULL,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "processed_avl_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "daily_reports" (
    "id" SERIAL NOT NULL,
    "device_id" INTEGER NOT NULL,
    "report_date" TIMESTAMP(3) NOT NULL,
    "total_data" INTEGER NOT NULL DEFAULT 0,
    "total_moving" INTEGER NOT NULL DEFAULT 0,
    "total_idle" INTEGER NOT NULL DEFAULT 0,
    "total_overspeed" INTEGER NOT NULL DEFAULT 0,
    "avg_speed" DOUBLE PRECISION NOT NULL DEFAULT 0,
    "max_speed" INTEGER NOT NULL DEFAULT 0,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "daily_reports_pkey" PRIMARY KEY ("id")
);

-- CreateIndex
CREATE UNIQUE INDEX "devices_imei_key" ON "devices"("imei");

-- CreateIndex
CREATE INDEX "processed_avl_device_id_idx" ON "processed_avl"("device_id");

-- CreateIndex
CREATE INDEX "processed_avl_timestamp_idx" ON "processed_avl"("timestamp");

-- CreateIndex
CREATE INDEX "daily_reports_device_id_idx" ON "daily_reports"("device_id");

-- CreateIndex
CREATE INDEX "daily_reports_report_date_idx" ON "daily_reports"("report_date");

-- AddForeignKey
ALTER TABLE "processed_avl" ADD CONSTRAINT "processed_avl_device_id_fkey" FOREIGN KEY ("device_id") REFERENCES "devices"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "daily_reports" ADD CONSTRAINT "daily_reports_device_id_fkey" FOREIGN KEY ("device_id") REFERENCES "devices"("id") ON DELETE RESTRICT ON UPDATE CASCADE;
