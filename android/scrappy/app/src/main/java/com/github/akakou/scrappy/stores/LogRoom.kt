package com.github.akakou.scrappy.stores

import androidx.room.*

@Entity(tableName = "signer_log")
data class SignerLog(
    @PrimaryKey(autoGenerate = true) val id: Int = 0,
    @ColumnInfo(name = "basename") val basename: String?,
)

@Dao
interface SignerLogDao {
    @Query("SELECT * FROM signer_log")
    fun getAll(): List<SignerLog>

    @Query("SELECT 1 FROM signer_log WHERE basename = :basename")
    fun hasExist(basename: String): Int

    @Insert
    fun insertAll(vararg signer_logs: SignerLog)

    @Delete
    fun delete(signer_log: SignerLog)
}

@Database(entities = [SignerLog::class], version = 2)
abstract class AppDatabase : RoomDatabase() {
    abstract fun signerLogDao(): SignerLogDao
}