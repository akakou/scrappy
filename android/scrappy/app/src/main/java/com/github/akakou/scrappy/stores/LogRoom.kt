package com.github.akakou.scrappy.stores

import androidx.room.*

@Entity(
    tableName = "signer_log",
    indices = [Index(value = ["i", "period", "basename"], unique = true)]
)
data class SignerLog(
    @PrimaryKey(autoGenerate = true) val id: Int = 0,
    @ColumnInfo(name = "i") val i: Int,
    @ColumnInfo(name = "period") val period: Long,
    @ColumnInfo(name = "basename") val basename: String,
)

@Dao
interface SignerLogDao {
    @Query("SELECT * FROM signer_log")
    fun getAll(): List<SignerLog>

    @Query("SELECT i FROM signer_log WHERE period = :period AND basename = :basename")
    fun getUsedIndexes(period: Long, basename: String): List<Int>

    @Insert
    fun insertAll(vararg signer_logs: SignerLog)

    @Delete
    fun delete(signer_log: SignerLog)
}

@Database(entities = [SignerLog::class], version = 3)
abstract class AppDatabase : RoomDatabase() {
    abstract fun signerLogDao(): SignerLogDao
}