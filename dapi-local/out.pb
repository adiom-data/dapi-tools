
�1
google/protobuf/timestamp.protogoogle.protobuf";
	Timestamp
seconds (Rseconds
nanos (RnanosB�
com.google.protobufBTimestampProtoPZ2google.golang.org/protobuf/types/known/timestamppb��GPB�Google.Protobuf.WellKnownTypesJ�/
 �
�
 2� Protocol Buffers - Google's data interchange format
 Copyright 2008 Google Inc.  All rights reserved.
 https://developers.google.com/protocol-buffers/

 Redistribution and use in source and binary forms, with or without
 modification, are permitted provided that the following conditions are
 met:

     * Redistributions of source code must retain the above copyright
 notice, this list of conditions and the following disclaimer.
     * Redistributions in binary form must reproduce the above
 copyright notice, this list of conditions and the following disclaimer
 in the documentation and/or other materials provided with the
 distribution.
     * Neither the name of Google Inc. nor the names of its
 contributors may be used to endorse or promote products derived from
 this software without specific prior written permission.

 THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
 "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
 LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
 A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
 OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
 SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
 LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
 DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
 THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
 (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
 OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.


  

" 
	
" 

# I
	
# I

$ ,
	
$ ,

% /
	
% /

& "
	

& "

' !
	
$' !

( ;
	
%( ;
�
 � �� A Timestamp represents a point in time independent of any time zone or local
 calendar, encoded as a count of seconds and fractions of seconds at
 nanosecond resolution. The count is relative to an epoch at UTC midnight on
 January 1, 1970, in the proleptic Gregorian calendar which extends the
 Gregorian calendar backwards to year one.

 All minutes are 60 seconds long. Leap seconds are "smeared" so that no leap
 second table is needed for interpretation, using a [24-hour linear
 smear](https://developers.google.com/time/smear).

 The range is from 0001-01-01T00:00:00Z to 9999-12-31T23:59:59.999999999Z. By
 restricting to that range, we ensure that we can convert to and from [RFC
 3339](https://www.ietf.org/rfc/rfc3339.txt) date strings.

 # Examples

 Example 1: Compute Timestamp from POSIX `time()`.

     Timestamp timestamp;
     timestamp.set_seconds(time(NULL));
     timestamp.set_nanos(0);

 Example 2: Compute Timestamp from POSIX `gettimeofday()`.

     struct timeval tv;
     gettimeofday(&tv, NULL);

     Timestamp timestamp;
     timestamp.set_seconds(tv.tv_sec);
     timestamp.set_nanos(tv.tv_usec * 1000);

 Example 3: Compute Timestamp from Win32 `GetSystemTimeAsFileTime()`.

     FILETIME ft;
     GetSystemTimeAsFileTime(&ft);
     UINT64 ticks = (((UINT64)ft.dwHighDateTime) << 32) | ft.dwLowDateTime;

     // A Windows tick is 100 nanoseconds. Windows epoch 1601-01-01T00:00:00Z
     // is 11644473600 seconds before Unix epoch 1970-01-01T00:00:00Z.
     Timestamp timestamp;
     timestamp.set_seconds((INT64) ((ticks / 10000000) - 11644473600LL));
     timestamp.set_nanos((INT32) ((ticks % 10000000) * 100));

 Example 4: Compute Timestamp from Java `System.currentTimeMillis()`.

     long millis = System.currentTimeMillis();

     Timestamp timestamp = Timestamp.newBuilder().setSeconds(millis / 1000)
         .setNanos((int) ((millis % 1000) * 1000000)).build();

 Example 5: Compute Timestamp from Java `Instant.now()`.

     Instant now = Instant.now();

     Timestamp timestamp =
         Timestamp.newBuilder().setSeconds(now.getEpochSecond())
             .setNanos(now.getNano()).build();

 Example 6: Compute Timestamp from current time in Python.

     timestamp = Timestamp()
     timestamp.GetCurrentTime()

 # JSON Mapping

 In JSON format, the Timestamp type is encoded as a string in the
 [RFC 3339](https://www.ietf.org/rfc/rfc3339.txt) format. That is, the
 format is "{year}-{month}-{day}T{hour}:{min}:{sec}[.{frac_sec}]Z"
 where {year} is always expressed using four digits while {month}, {day},
 {hour}, {min}, and {sec} are zero-padded to two digits each. The fractional
 seconds, which can go up to 9 digits (i.e. up to 1 nanosecond resolution),
 are optional. The "Z" suffix indicates the timezone ("UTC"); the timezone
 is required. A proto3 JSON serializer should always use UTC (as indicated by
 "Z") when printing the Timestamp type and a proto3 JSON parser should be
 able to accept both UTC and other timezones (as indicated by an offset).

 For example, "2017-01-15T01:30:15.01Z" encodes 15.01 seconds past
 01:30 UTC on January 15, 2017.

 In JavaScript, one can convert a Date object to this format using the
 standard
 [toISOString()](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Date/toISOString)
 method. In Python, a standard `datetime.datetime` object can be converted
 to this format using
 [`strftime`](https://docs.python.org/2/library/time.html#time.strftime) with
 the time format spec '%Y-%m-%dT%H:%M:%S.%fZ'. Likewise, in Java, one can use
 the Joda Time's [`ISODateTimeFormat.dateTime()`](
 http://joda-time.sourceforge.net/apidocs/org/joda/time/format/ISODateTimeFormat.html#dateTime()
 ) to obtain a formatter capable of generating timestamps in this format.



 �
�
  �� Represents seconds of UTC time since Unix epoch
 1970-01-01T00:00:00Z. Must be from 0001-01-01T00:00:00Z to
 9999-12-31T23:59:59Z inclusive.


  �

  �

  �
�
 �� Non-negative fractions of a second at nanosecond resolution. Negative
 second values with fractions must still have non-negative nanos values
 that count forward in time. Must be from 0 to 999,999,999
 inclusive.


 �

 �

 �bproto3�� 
�4
	out.protocom.examplegoogle/protobuf/timestamp.proto"=
MongoBinary
_subtype (RSubtype
_data (RData"h
TestMoviesTomatoesViewer
meter (Rmeter

numReviews (R
numReviews
rating (Rrating"h
TestMoviesTomatoesCritic
meter (Rmeter

numReviews (R
numReviews
rating (Rrating"�
TestMoviesTomatoes
	boxOffice (	R	boxOffice
	consensus (	R	consensus=
critic (2%.com.example.TestMoviesTomatoesCriticRcritic,
dvd (2.google.protobuf.TimestampRdvd
fresh (Rfresh<
lastUpdated (2.google.protobuf.TimestampRlastUpdated

production (	R
production
rotten (Rrotten=
viewer	 (2%.com.example.TestMoviesTomatoesViewerRviewer
website
 (	Rwebsite"N
TestMoviesImdb
id (Rid
rating (Rrating
votes (Rvotes"\
TestMoviesAwards 
nominations (Rnominations
text (	Rtext
wins (Rwins"�

TestMovies
_id (	RId5
awards (2.com.example.TestMoviesAwardsRawards
cast (	Rcast
	countries (	R	countries
	directors (	R	directors
fullplot (	Rfullplot
genres (	Rgenres/
imdb (2.com.example.TestMoviesImdbRimdb
	languages	 (	R	languages 
lastupdated
 (	Rlastupdated

metacritic (R
metacritic,
num_mflix_comments (RnumMflixComments
plot (	Rplot
poster (	Rposter
rated (	Rrated6
released (2.google.protobuf.TimestampRreleased
runtime (Rruntime
title (	Rtitle;
tomatoes (2.com.example.TestMoviesTomatoesRtomatoes
type (	Rtype
writers (	Rwriters
year (Ryear"E
ListTestMoviesRequest
cursor (	Rcursor
limit (Rlimit"|
ListTestMoviesResponse+
data (2.com.example.TestMoviesRdata
next_cursor (	R
nextCursor
limit (Rlimit"'
GetTestMoviesRequest
_id (	RId"D
GetTestMoviesResponse+
data (2.com.example.TestMoviesRdata"F
CreateTestMoviesRequest+
data (2.com.example.TestMoviesRdata"+
CreateTestMoviesResponse
_id (	RId"F
UpdateTestMoviesRequest+
data (2.com.example.TestMoviesRdata"
UpdateTestMoviesResponse"*
DeleteTestMoviesRequest
_id (	RId"
DeleteTestMoviesResponse2�
ExampleServiceY
ListTestMovies".com.example.ListTestMoviesRequest#.com.example.ListTestMoviesResponseV
GetTestMovies!.com.example.GetTestMoviesRequest".com.example.GetTestMoviesResponse_
CreateTestMovies$.com.example.CreateTestMoviesRequest%.com.example.CreateTestMoviesResponse_
UpdateTestMovies$.com.example.UpdateTestMoviesRequest%.com.example.UpdateTestMoviesResponse_
DeleteTestMovies$.com.example.DeleteTestMoviesRequest%.com.example.DeleteTestMoviesResponseJ�
  v#

  

 
	
  )


  	


 

  

  

  

  

 

 

 

 


 


 

 

 

 

 













	




 


 

 

 

 

 













	




 "




 

 

 	

 





	



&



!

$%

$





"#









,



'

*+





	











 &

 

 !

 $%

	!

	!

	!	

	!


$ (


$

 %

 %

 %


 %

&

&

&	

&

'

'

'

'


* .


*

 +

 +

 +

 +

,

,

,	

,

-

-

-

-


0 G


0

 1

 1

 1	

 1

2

2

2

2

3

3


3

3

3

4 

4


4

4

4

5 

5


5

5

5

6

6

6	

6

7

7


7

7

7

8

8

8

8

9 

9


9

9

9

	:

	:

	:	

	:


;


;


;


;

< 

<

<

<

=

=

=	

=

>

>

>	

>

?

?

?	

?

@*

@

@$

@')

A

A

A

A

B

B

B	

B

C#

C

C

C "

D

D

D	

D

E

E


E

E

E

F

F

F

F


 I O


 I

  JM

  J

  J*

  J5K

 KJ

 K

 K(

 K3H

 LS

 L

 L.

 L9Q

 MS

 M

 M.

 M9Q

 NS

 N

 N.

 N9Q


Q T


Q

 R

 R

 R	

 R

S

S

S

S


V Z


V

 W

 W


 W

 W

 W

X

X

X	

X

Y

Y

Y

Y


	\ ^


	\

	 ]

	 ]

	 ]	

	 ]



` b



`


 a


 a


 a


 a


d f


d

 e

 e

 e

 e


h j


h 

 i

 i

 i	

 i


l n


l

 m

 m

 m

 m
	
p #


p 


r t


r

 s

 s

 s	

 s
	
v #


v bproto3��  