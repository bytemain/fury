/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

import { alloc } from '@furyjs/fury/lib/platformBuffer';
import { BinaryReader } from '@furyjs/fury/lib/reader';
import { Config } from '@furyjs/fury/lib/type';
import { BinaryWriter } from '@furyjs/fury/lib/writer';
import { describe, expect, test } from '@jest/globals';
const hps = process.env.enableHps ? require('@furyjs/hps') : null;



[
    {
        useLatin1: true,
        hps,
    },
    {
        useLatin1: false,
        hps,
    }
].forEach((config: Config) => {
    describe('writer', () => {
        test('should uint8 work', () => {
            const writer = BinaryWriter(config);
            {
                writer.uint8(10);
                var ab = writer.dump();
                expect(ab.byteLength).toBe(1);
                expect(ab[0]).toBe(10);
                expect(writer.getCursor()).toBe(1);
            }

            {
                writer.uint8(256);
                var ab = writer.dump();

                expect(ab.byteLength).toBe(2);
                expect(ab[1]).toBe(0);
                expect(writer.getCursor()).toBe(2);
            }
        });
    });
})


describe('reader', () => {

    test('should bufferRef work', () => {
        const bb = alloc(100);
        bb.latin1Write("hello", 0);
        const target = new Uint8Array(5);
        bb.copy(target, 0, 0, 5);
        expect([...target]).toEqual([ 104, 101, 108, 108, 111 ])

        const reader = BinaryReader({});

        reader.reset(bb);
        const ref = reader.bufferRef(5);
        ref[0] = 0;
        bb.copy(target, 0, 0, 5);
        expect([...target]).toEqual([ 0, 101, 108, 108, 111 ])
    })
})
