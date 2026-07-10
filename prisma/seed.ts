import "dotenv/config";
import { PrismaClient } from "../generated/prisma/client.js";
import { PrismaPg } from "@prisma/adapter-pg";

const adapter = new PrismaPg({
  connectionString: process.env.DATABASE_URL!,
});

const prisma = new PrismaClient({ adapter });

async function main() {
  await prisma.device.upsert({
    where: {
      imei: "354002390725990",
    },
    update: {
      deviceName: "Teltonika FMB920 Company",
      isActive: true,
    },
    create: {
      imei: "354002390725990",
      deviceName: "Teltonika FMB920 Company",
      isActive: true,
    },
  });

  console.log("Seed device berhasil");
}

main()
  .catch((error) => {
    console.error(error);
    process.exit(1);
  })
  .finally(async () => {
    await prisma.$disconnect();
  });