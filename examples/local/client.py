#!/usr/bin/env python

# Copyright 2017 Google Inc.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
"""Sample Vitess client in Python.

This is a sample for using the Python Vitess client.
It's a script that inserts some random messages on random pages of the
guestbook sample app.

Before running this, start up a local example cluster as described in the
README.md file.

Then run client.sh, which sets up PYTHONPATH before running client.py:
vitess/examples/local$ ./client.sh
"""

import argparse
import random
import time

from vtdb import vtgate_client

# register the python gRPC client upon import
from vtdb import grpc_vtgate_client  # pylint: disable=unused-import

# Parse args
parser = argparse.ArgumentParser()
parser.add_argument('--server', dest='server', default='localhost:15991')
parser.add_argument('--timeout', dest='timeout', type=float, default='10.0')
args = parser.parse_args()

# Connect
conn = vtgate_client.connect('grpc', args.server, args.timeout)

try:
    # Insert some messages on random pages.
    print 'Inserting into bypage.messages'
    cursor = conn.cursor(tablet_type='master', writable=True)

    cursor.begin()
    for i in range(10000):
        page = random.randint(1, 100)

        cursor.execute(
            'INSERT INTO bypage.messages (page, time_created_ns, message)'
            ' VALUES (:page, :time_created_ns, :message)', {
                'page': page,
                'time_created_ns': int(time.time() * 1e9),
                'message': 'V is for speed',
            })
    cursor.commit()

    print 'Inserting into byuserid.users'
    cursor.begin()
    for i in range(10000):
        cursor.execute(
            'INSERT INTO byuserid.users (userid, name, email)'
            ' VALUES (:userid, :name, :email)', {
                'userid': i,
                'name': "test user",
                "email": "test@test.com"
            })
    cursor.commit()

    print 'Inserting into byuserid.files'
    cursor.begin()
    for i in range(10000):
        cursor.execute(
            'INSERT INTO byuserid.files (fileid, userid, content)'
            ' VALUES (:fileid, :userid, :content)', {
                'fileid': i,
                'userid': i,
                "content": "file content"
            })
    cursor.commit()

    # Read it back from the master.
    print 'Reading from master...'

    cursor.execute(
        'SELECT page, time_created_ns, message FROM messages LIMIT 10', {})
    for row in cursor.fetchall():
        print row

    cursor.execute('SELECT userid, name, email FROM users LIMIT 10', {})
    for row in cursor.fetchall():
        print row

    cursor.close()

    # Read from a replica.
    # Note that this may be behind master due to replication lag.
    print 'Reading from replica...'
    cursor = conn.cursor(tablet_type='replica')

    cursor.execute('SELECT page, time_created_ns, message FROM messages LIMIT 10', {})
    for row in cursor.fetchall():
        print row

    cursor.execute('SELECT userid, name, email FROM users LIMIT 10', {})
    for row in cursor.fetchall():
        print row

    cursor.close()

finally:
    # Clean up
    conn.close()
