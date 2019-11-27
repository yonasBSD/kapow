Sharing the Stats
=================

**Junior**

  Good morning!

**Senior**

  Just about time...  We are in trouble!

  The report stuff was a complete success, so much that now *Susan* has
  hired a frontend developer to create a custom dashboard to see the
  stats in real time.

  Now we have to provide the backend for the solution.

**Junior**

  And whats the problem?

**Senior**

  We are not developers what are we doing writing backend?

**Junior**

  Chill out, man. Can't be that difficult?  What they need exactly?

**Senior**

  We have to provide a new endpoint to serve the same data but in JSON
  format.

**Junior**

  So we have half of the work already done.

  What about this?

  .. code-block:: bash

     kapow route add /capacitystats - <<-HERE 
       echo "{\"memory\": \"`free -m`\"}"  | kapow set /response/body
     HERE

**Senior**

  For starters that's not valid JSON. The output would be something
  like:

  .. code-block:: console

     $ echo "{\"memory\": \"`free -m`\"}"
     {"memory": "              total        used        free      shared  buff/cache   available
     Mem:          31967        3121       21680         980        7166       27418
     Swap:             0           0           0"}

  You can't add new lines inside a JSON string that way, you have to
  encode with ``\n``.


**Junior**

  Are you sure?

**Senior**

  See it by yourself.

  .. code-block:: console

     $ echo "{\"memory\": \"`free -m`\"}" | jq
     parse error: Invalid string: control characters from U+0000 through U+001F must be escaped at line 3, column 44

**Junior**

  ``jq``? What is that command?

**Senior**

  ``jq`` is a wonderful tool for working with JSON data from the command
  line.  With you ``jq`` you can extract data from JSON and also
  generate well-formed JSON.

**Junior**

  Let's use it then! 

  How can we generate a JSON document with ``jq``?

**Senior**

  To generate a document we use the ``-n`` argument:

  .. code-block:: console

     $ jq -n '{"mykey": "myvalue"}'
     {
       "mykey": "myvalue"
     }

**Junior**

  That is not very useful. The output is the same.

**Senior**

  It get's better. You can add variables to the JSON and ``jq`` will escape them for you.

  .. code-block:: console

     $ jq -n --arg myvar "$(echo -n myvalue)" '{"mykey": $myvar}'
     {
       "mykey": "myvalue"
     }

**Junior**

  That's just what I need.

  What do you think of this?

  .. code-block:: console

     $ jq -n --arg host "$(hostname)" --arg date "$(date)" --arg memory "$(free -m)" --arg load "$(uptime)" --arg disk "$(df -h)" '{"hostname": $host, "date": $date, "memory": $memory, "load": $load, "disk": $disk}'
     {
       "hostname": "junior-host",
       "date": "Tue 26 Nov 2019 05:27:24 PM CET",
       "memory": "              total        used        free      shared  buff/cache   available\nMem:          31967        3114       21744         913        7109       27492\nSwap:             0           0           0",
       "load": " 17:27:24 up 10:21,  1 user,  load average: 0.20, 0.26, 0.27",
       "disk": "Filesystem          Size  Used Avail Use% Mounted on\ndev                  16G     0   16G   0% /dev"
     }

**Senior**

  That is the data we have to produce.  But the code is far from readable.  And
  you also forgot about adding the endpoint.

  Can we do any better?

**Junior**

  That's easy:

  .. code-block:: bash

     kapow route add /capacitystats - <<-HERE 
       jq -n \
          --arg hostname "$(hostname)" \
          --arg date "$(date)" \
          --arg memory "$(free -m)" \
          --arg load "$(uptime)" \
          --arg disk "$(df -h)" \
          '{"hostname": $hostname, "date": $date, "memory": $memory, "load": $load, "disk": $disk}' \
       | kapow set /response/body
     HERE

  What do you think?
   
**Senior**

  You forgot one more thing.

**Junior**

  I think you are wrong, the JSON is well-formed and it contains all the
  required data.  Also the code is very readable.

**Senior**

  You are right but, you are not using HTTP correctly.  You have to set the
  ``Content-Type`` header to let your client know the format of the data you are
  outputting.

**Junior**

  Ok, let me try:

  .. code-block:: bash

     kapow route add /capacitystats - <<-HERE 
       jq -n \
          --arg hostname "$(hostname)" \
          --arg date "$(date)" \
          --arg memory "$(free -m)" \
          --arg load "$(uptime)" \
          --arg disk "$(df -h)" \
          '{"hostname": $hostname, "date": $date, "memory": $memory, "load": $load, "disk": $disk}' \
       | kapow set /response/body
       echo application/json | kapow set /response/headers/Content-Type
     HERE

**Senior**

  Just a couple of details.

  1. You have to set the headers **before** the body.  This is because the body
     can be so big that Kapow! is forced to start sending it out.
  2. In cases where you want to set a small piece of data (like the header) is
     better to not use the ``stdin``.  Kapow! provides a secondary syntax for these
     cases:

     .. code-block:: console

        $ kapow set <resource> <value>

**Junior**

  Something like this?

  .. code-block:: bash

     kapow route add /capacitystats - <<-HERE 
       kapow set /response/headers/Content-Type application/json
       jq -n \
          --arg hostname "$(hostname)" \
          --arg date "$(date)" \
          --arg memory "$(free -m)" \
          --arg load "$(uptime)" \
          --arg disk "$(df -h)" \
          '{"hostname": $hostname, "date": $date, "memory": $memory, "load": $load, "disk": $disk}' \
       | kapow set /response/body
     HERE

**Senior**

  That's perfect!  Let's upload this to the *Corporate Server* and tell the
  frontend developer.