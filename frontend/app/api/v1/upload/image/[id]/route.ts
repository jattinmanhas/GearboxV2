import { NextRequest, NextResponse } from "next/server";
import { unlink } from "fs/promises";
import { join } from "path";
import { isCloudinaryConfigured } from "@/lib/image-upload";

const UPLOAD_DIR = join(process.cwd(), "public", "uploads");

export async function DELETE(
  request: NextRequest,
  { params }: { params: Promise<{ id: string }> }
) {
  try {
    // Await params to get the id
    const { id } = await params;

    // Decode the URL-encoded public ID (handles slashes like gearbox-uploads/filename)
    const decodedId = decodeURIComponent(id);

    const body = await request.json();
    const isCloudinary = body?.isCloudinary || false;

    if (isCloudinary && isCloudinaryConfigured()) {
      // Delete from Cloudinary
      try {
        const { deleteFromCloudinary } = await import("@/lib/cloudinary");
        await deleteFromCloudinary(decodedId);

        return NextResponse.json({
          success: true,
          message: "Image deleted from Cloudinary successfully",
        });
      } catch (error) {
        console.error("Cloudinary delete error:", error);
        return NextResponse.json(
          {
            success: false,
            message: "Failed to delete from Cloudinary",
            error: error instanceof Error ? error.message : "Unknown error",
          },
          { status: 500 }
        );
      }
    } else {
      // Delete local file
      const filename = decodedId;
      const filePath = join(UPLOAD_DIR, filename);

      try {
        await unlink(filePath);

        // Also delete thumbnails if they exist
        const thumbnailDir = join(UPLOAD_DIR, "thumbnails");
        const thumbnailFiles = [
          `${filename.split(".")[0]}_thumbnail.jpg`,
          `${filename.split(".")[0]}_small.jpg`,
          `${filename.split(".")[0]}_medium.jpg`,
          `${filename.split(".")[0]}_large.jpg`,
        ];

        for (const thumbnailFile of thumbnailFiles) {
          try {
            await unlink(join(thumbnailDir, thumbnailFile));
          } catch {
            // Thumbnail might not exist, ignore error
          }
        }

        return NextResponse.json({
          success: true,
          message: "Image deleted successfully",
        });
      } catch (error) {
        console.error("Local delete error:", error);
        return NextResponse.json(
          {
            success: false,
            message: "Failed to delete image",
            error: error instanceof Error ? error.message : "Unknown error",
          },
          { status: 500 }
        );
      }
    }
  } catch (error) {
    console.error("Delete error:", error);
    return NextResponse.json(
      {
        success: false,
        message: "Delete failed",
        error: error instanceof Error ? error.message : "Unknown error",
      },
      { status: 500 }
    );
  }
}