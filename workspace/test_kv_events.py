import json
import sys
from typing import Any

import msgspec
import zmq
from msgspec.msgpack import Decoder

from sglang.srt.disaggregation.kv_events import (
    AllBlocksCleared,
    BlockRemoved,
    BlockStored,
    KVEventBatch,
)

def format_event(event: Any) -> dict:
    """Format an event for display."""
    if isinstance(event, BlockStored):
        return {
            "type": "BlockStored",
            "block_hashes": event.block_hashes,
            "parent_block_hash": event.parent_block_hash,
            # "token_ids": event.token_ids[:10] if len(event.token_ids) > 10 else event.token_ids,  # Show first 10
            "token_ids": event.token_ids, # if len(event.token_ids) > 10 else event.token_ids,  # Show first 10
            "block_size": event.block_size,
            "lora_id": event.lora_id,
        }
    elif isinstance(event, BlockRemoved):
        return {
            "type": "BlockRemoved",
            "block_hashes": event.block_hashes,
        }
    elif isinstance(event, AllBlocksCleared):
        return {
            "type": "AllBlocksCleared",
        }
    else:
        return {"type": str(type(event).__name__), "data": str(event)}

if __name__ == '__main__':
    decoder = Decoder(type=KVEventBatch)
    context = zmq.Context()
    sub = context.socket(zmq.SUB)

    #sub.connect("tcp://localhost:5557")
    sub.bind("tcp://*:5557")

    sub.setsockopt_string(zmq.SUBSCRIBE, "")

    try:
        while True:
            # Receive multipart message: (topic, sequence, payload)
            try:
                topic_bytes, seq_bytes, payload = sub.recv_multipart()
                seq = int.from_bytes(seq_bytes, "big")

                # Decode the event batch
                event_batch = decoder.decode(payload)

                # Process each event in the batch
                for event in event_batch.events:
                    event_data = format_event(event)

                    if True:
                        output = {
                            "sequence": seq,
                            "timestamp": event_batch.ts,
                            "attn_dp_rank": event_batch.attn_dp_rank,
                            "event": event_data,
                        }

                        print(json.dumps(output, indent=2))
                    else:
                        print(f"[Seq {seq}] [TS {event_batch.ts:.3f}] [DP Rank {event_batch.attn_dp_rank}]")
                        print(f"  {json.dumps(event_data, indent=2)}")
                        print()

            except KeyboardInterrupt:
                print("\nStopping...")
                break
            except Exception as e:
                print(f"Error receiving/decoding event: {e}", file=sys.stderr)
                continue

    finally:
        sub.close()
        context.term()
