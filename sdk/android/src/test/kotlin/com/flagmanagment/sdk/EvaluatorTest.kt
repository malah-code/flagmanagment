package com.flagmanagment.sdk

import org.junit.Assert.assertEquals
import org.junit.Assert.assertTrue
import org.junit.Test

class EvaluatorTest {

    @Test
    fun testBucketUserRange() {
        for (i in 0 until 1000) {
            val bucket = Evaluator.bucketUser("user-$i")
            assertTrue("Bucket must be in [0, 99]", bucket in 0..99)
        }
    }

    @Test
    fun testBucketUserDeterministic() {
        val b1 = Evaluator.bucketUser("user_consistent_target")
        val b2 = Evaluator.bucketUser("user_consistent_target")
        assertEquals("Same targeting key must produce identical bucket", b1, b2)
    }
}
