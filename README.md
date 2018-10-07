# Zeats

Zeats service maintains ice-cream details by providing CRUD functionality.

Implemented a RESTful API with CRUD functionality.
Written in Golang, adhering to Golang conventions and best practices.
Uses cassandra databse to store ice-cream details.
Secured the API with JSON Web Tokens (JWT) authentication mechanism.

## Reads
* Golang all abouts: https://golang.org/doc/install
* Golang code style: https://godoc.org/golang.org/x/tools/cmd/goimports
* Glide: https://github.com/Masterminds/glide
* Cassandra: https://wiki.apache.org/cassandra

## Setup
* Install cassandra 2.2.10 or above
* Start cassandra server using below command
```$xslt
   /path/to/cassandra/apache-cassandra-2.2.10/bin/cassandra 
```
* Create keyspace `zeats` using command below
```$xslt
CREATE KEYSPACE zeats WITH replication = {'class': 'SimpleStrategy', 'replication_factor': '1'}  AND 
durable_writes = true;
```
* Create set eats 
```$xslt
CREATE TABLE zeats.eats (
    product_id text PRIMARY KEY,
    allergy_info text,
    description text,
    dietary_certifications text,
    image_closed text,
    image_open text,
    ingredients list<text>,
    name text,
    sourcing_values list<text>,
    story text
) WITH bloom_filter_fp_chance = 0.01
    AND caching = '{"keys":"ALL", "rows_per_partition":"NONE"}'
    AND comment = ''
    AND compaction = {'class': 'org.apache.cassandra.db.compaction.SizeTieredCompactionStrategy'}
    AND compression = {'sstable_compression': 'org.apache.cassandra.io.compress.LZ4Compressor'}
    AND dclocal_read_repair_chance = 0.1
    AND default_time_to_live = 0
    AND gc_grace_seconds = 864000
    AND max_index_interval = 2048
    AND memtable_flush_period_in_ms = 0
    AND min_index_interval = 128
    AND read_repair_chance = 0.0
    AND speculative_retry = '99.0PERCENTILE';
```
* Install go1.10.3 or above
* Create a workspace for go and create `src`, `pkg`, `logs` and `bin` directory inside it
* Install zeats
```$xslt
cd /path/to/your/workspace/src
git clone git@github.com:henzilfernandes/zeats.git
cd zeats
glide install -v .
go install -v .
../../bin/zeats -log_dir=../../logs -v=0
```

### APIs

* Create Authentication Token

Request :
```$xslt
curl -X POST \
  http://localhost:8888/v1/createToken \
  -H 'accept: application/json' \
  -H 'cache-control: no-cache' \
  -H 'content-type: application/json' \
  -d '{
	"username": "your_username",
	"password": "your_password"
}'

```

Response :
```$xslt
{
    "statusCode": 100,
    "statusMessage": "Authentication Token created successfully",
    "token": "your_private_token"
}
```

* Insert Eats

Request :
```$xslt
curl -X POST \
  http://localhost:8888/v1/insert \
  -H 'accept: application/json' \
  -H 'authorization: your_private_token' \
  -H 'cache-control: no-cache' \
  -H 'content-type: application/json' \
  -d '{
  "name": "Vanilla Toffee Bar Crunch",
  "image_closed": "/files/live/sites/systemsite/files/flavors/products/us/pint/open-closed-pints/vanilla-toffee-landing.png",
  "image_open": "/files/live/sites/systemsite/files/flavors/products/us/pint/open-closed-pints/vanilla-toffee-landing-open.png",
  "description": "Vanilla Ice Cream with Fudge-Covered Toffee Pieces",
  "story": "Vanilla What Bar Crunch? We gave this flavor a new name to go with the new toffee bars we’re using as part of our commitment to source Fairtrade Certified and non-GMO ingredients. We love it and know you will too!",
  "sourcing_values": [
    "Non-GMO",
    "Cage-Free Eggs",
    "Fairtrade",
    "Responsibly Sourced Packaging",
    "Caring Dairy"
  ],
  "ingredients": [
    "cream",
    "skim milk",
    "liquid sugar",
    "water",
    "sugar",
    "coconut oil",
    "egg yolks",
    "butter",
    "vanilla extract",
    "almonds",
    "cocoa (processed with alkali)",
    "milk",
    "soy lecithin",
    "cocoa",
    "natural flavor",
    "salt",
    "vegetable oil",
    "guar gum",
    "carrageenan"
  ],
  "allergy_info": "may contain wheat, peanuts and other tree nuts",
  "dietary_certifications": "Kosher",
  "productId": "646"
}'
```
Response :
```$xslt
{
    "statusCode": 100,
    "statusMessage": "Product data inserted successfully"
}
```
* Fetch Eats

Request :
```$xslt
curl -X GET \
  http://localhost:8888/v1/fetch/646 \
  -H 'accept: application/json' \
  -H 'authorization: your_private_token' \
  -H 'cache-control: no-cache' \
  -H 'content-type: application/json' 
```
Response :
```$xslt
{
    "statusCode": 100,
    "statusMessage": "Product data fetched successfully",
    "payload": {
        "productId": "646",
        "name": "Vanilla Toffee Bar Crunch",
        "image_closed": "/files/live/sites/systemsite/files/flavors/products/us/pint/open-closed-pints/vanilla-toffee-landing.png",
        "image_open": "/files/live/sites/systemsite/files/flavors/products/us/pint/open-closed-pints/vanilla-toffee-landing-open.png",
        "description": "Vanilla Ice Cream with Fudge-Covered Toffee Pieces",
        "story": "Vanilla What Bar Crunch? We gave this flavor a new name to go with the new toffee bars we’re using as part of our commitment to source Fairtrade Certified and non-GMO ingredients. We love it and know you will too!",
        "sourcing_values": [
            "Non-GMO",
            "Cage-Free Eggs",
            "Fairtrade",
            "Responsibly Sourced Packaging",
            "Caring Dairy"
        ],
        "ingredients": [
            "cream",
            "skim milk",
            "liquid sugar",
            "water",
            "sugar",
            "coconut oil",
            "egg yolks",
            "butter",
            "vanilla extract",
            "almonds",
            "cocoa (processed with alkali)",
            "milk",
            "soy lecithin",
            "cocoa",
            "natural flavor",
            "salt",
            "vegetable oil",
            "guar gum",
            "carrageenan"
        ],
        "allergy_info": "may contain wheat, peanuts and other tree nuts",
        "dietary_certifications": "Kosher"
    }
}
```
* Update Eats

Request :
```$xslt
curl -X POST \
  http://localhost:8888/v1/update \
  -H 'accept: application/json' \
  -H 'authorization: your_private_token' \
  -H 'cache-control: no-cache' \
  -H 'content-type: application/json' \
  -d '{
  "description": "Vanilla Ice Cream with Fudge-Covered Toffee Pieces .......",
  "sourcing_values": [
    "Non-GMO",
    "Responsibly Sourced Packaging",
    "Caring Dairy"
  ],
  "ingredients": [
    "cream",
    "skim milk",
    "liquid sugar",
    "water",
    "sugar",
    "coconut oil",
    "egg yolks",
    "butter",
    "vanilla extract",
    "almonds",
    "milk",
    "soy lecithin",
    "cocoa",
    "natural flavor",
    "salt",
    "vegetable oil",
    "guar gum",
    "carrageenan"
  ],
  "allergy_info": "may contain wheat, peanuts and other tree nuts",
  "dietary_certifications": "Putin",
  "productId": "646"
}'
```
Response :
```$xslt
{
    "statusCode": 100,
    "statusMessage": "Product data updated successfully"
}
```

* Delete Eats

Request :
```$xslt
curl -X DELETE \
  http://localhost:8888/v1/delete/646 \
  -H 'accept: application/json' \
  -H 'authorization: your_private_token' \
  -H 'cache-control: no-cache' \
  -H 'content-type: application/json' 
```

