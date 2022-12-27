import random
from confluent_kafka import avro
from confluent_kafka.avro import AvroProducer

schema_v1 = """
{
 "namespace": "io.confluent.examples.clients.basicavro",
 "type": "record",
 "name": "Payment",
 "fields": [
     {"name": "id", "type": "string"},
     {"name": "amount", "type": "double"}
 ]
}
"""

schema_v2 = """
{
 "namespace": "io.confluent.examples.clients.basicavro",
 "type": "record",
 "name": "Payment",
 "fields": [
     {"name": "id", "type": "string"},
     {"name": "amount", "type": "double"},
     {"name": "name", "type": "string"}
 ]
}
"""

value_schema = avro.loads(schema_v2)

producer_conf = {'bootstrap.servers': 'localhost:9092',
                 'schema.registry.url': 'http://localhost:8081'}

avroProducer = AvroProducer(producer_conf, default_value_schema=value_schema)

for i in range(1, 20):
    avroProducer.produce(topic='payments-topic',
                         value={"id": "transact_%s" % i,
                                "amount": random.uniform(10, 500),
                                "name": "customer_%s" % i,})
avroProducer.flush()
