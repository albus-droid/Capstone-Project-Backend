package listing

import (
    "net/http"
    "os"
    "fmt"
    "time"
    "github.com/gin-gonic/gin"
   	"github.com/albus-droid/Capstone-Project-Backend/internal/auth"
    "github.com/minio/minio-go/v7"
)


// RegisterRoutes mounts listing endpoints under /listings
func RegisterRoutes(r *gin.Engine, svc Service, minioClient *minio.Client) {
    grp := r.Group("/listings")
    grp.Use(auth.Middleware())

    // POST /listings
    grp.POST("", func(c *gin.Context) {
        var l Listing
        if err := c.ShouldBindJSON(&l); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }
        if err := svc.Create(&l); err != nil {
            c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
            return
        }
        c.JSON(http.StatusCreated, gin.H{"message": "listing created", "id": l.ID,})
    })

    // GET /listings/:id
    grp.GET("/:id", func(c *gin.Context) {
        id := c.Param("id")
        l, err := svc.GetByID(id)
        if err != nil {
            c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
            return
        }
        c.JSON(http.StatusOK, l)
    })

    // GET /listings?sellerId=...
    grp.GET("", func(c *gin.Context) {
        if sellerID := c.Query("sellerId"); sellerID != "" {
            list, _ := svc.ListBySeller(sellerID)
            c.JSON(http.StatusOK, list)
            return
        }
        // fallback: list all
        c.JSON(http.StatusOK, svc.ListAll())
    })

    // PUT /listings/:id  — update an existing listing
    grp.PUT("/:id", func(c *gin.Context) {
        id := c.Param("id")
        var l Listing
        if err := c.ShouldBindJSON(&l); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }
        if err := svc.Update(id, l); err != nil {
            // you can customize error handling based on your svc.Update error
            c.JSON(http.StatusNotFound, gin.H{"error": "not found or unable to update"})
            return
        }
        c.JSON(http.StatusOK, gin.H{"message": "listing updated"})
    })

    // DELETE /listings/:id  — remove a listing
    grp.DELETE("/:id", func(c *gin.Context) {
        id := c.Param("id")
        if err := svc.Delete(id); err != nil {
            c.JSON(http.StatusNotFound, gin.H{"error": "not found or unable to delete"})
            return
        }
        c.Status(http.StatusNoContent)
    })


    // POST /listings/:id/image — Uploads image to MinIO and saves URL to Listing in DB
    grp.POST("/:id/image", func(c *gin.Context) {
        listingID := c.Param("id")
        file, header, err := c.Request.FormFile("file")
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "no file uploaded"})
            return
        }
        defer file.Close()

        bucket := os.Getenv("MINIO_BUCKET")
        if bucket == "" {
            bucket = "listing-images"
        }
        objectName := fmt.Sprintf("listings/%s/%s", listingID, header.Filename)
        contentType := header.Header.Get("Content-Type")

        _, err = minioClient.PutObject(
            c, bucket, objectName, file, header.Size,
            minio.PutObjectOptions{ContentType: contentType},
        )
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "minio upload failed"})
            return
        }

        imageAPIUrl := fmt.Sprintf("/listings/%s/image/%s", listingID, header.Filename)

        // --- Update listing in DB ---
        listing, err := svc.GetByID(listingID)
        if err != nil {
            c.JSON(http.StatusNotFound, gin.H{"error": "listing not found"})
            return
        }

        // For single image per listing:
        listing.Image = imageAPIUrl

        // Save updated listing
        err = svc.Update(listingID, *listing)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update listing with image"})
            return
        }

        c.JSON(http.StatusOK, gin.H{
            "image_url": imageAPIUrl,
        })
    })

    // GET /listings/:id/image/:filename — Returns signed URL for image view
    grp.GET("/:id/image/:filename", func(c *gin.Context) {
        listingID := c.Param("id")
        filename := c.Param("filename")
        bucket := os.Getenv("MINIO_BUCKET")
        if bucket == "" {
            bucket = "listing-images"
        }
        objectName := fmt.Sprintf("listings/%s/%s", listingID, filename)

        signedURL, err := minioClient.PresignedGetObject(
            c, bucket, objectName, time.Hour, nil,
        )
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate signed URL"})
            return
        }

        c.JSON(http.StatusOK, gin.H{
            "signed_url": signedURL.String(),
        })
    })
}