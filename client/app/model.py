from sqlalchemy.ext.declarative import declarative_base
from sqlalchemy import create_engine, Column, Integer, Text, Float, DateTime
from sqlalchemy.orm import scoped_session, sessionmaker
from datetime import datetime
import random, sys, logging

sys.path.append('/attestation/common/')
from wrapper import TERM_LEN, MAX_COUNTER

handler = logging.FileHandler('/attestation/ignored_workspace/client.log')
handler.setLevel(logging.ERROR)
logging.getLogger('sqlalchemy').addHandler(handler)


MODE_ATTACK = False

ENGINE = create_engine('sqlite:////attestation/ignored_workspace/db.sqlite3', echo=False)
Base = declarative_base()

session = scoped_session(
    sessionmaker(
        bind = ENGINE,
        autocommit=False,
    )
)

ERRORS = {
    "HAVE_USED_ALL_COUNTER": 1
}

class AttestationLog(Base):
    __tablename__ = 'attestation_logs'

    _id = Column('id', Integer, primary_key = True)
    origin = Column('origin', Text)
    counter = Column('counter', Text)
    created_at = Column(DateTime, default=datetime.now)

    def find_by_origin_and_datetime(origin, datetime):
        attestation_logs = session.query(AttestationLog) \
                .filter(AttestationLog.origin == origin, AttestationLog.created_at == datetime)

        return attestation_logs

    def select_counters_by_origin_and_datetime(origin, datetime):
        attestation_logs = AttestationLog.find_by_origin_and_datetime(origin, datetime)
        return list(map(lambda x: x.origin, attestation_logs))

    def gen_new_counter(origin, now):
        counters = AttestationLog.select_counters_by_origin_and_datetime(origin, now) 
        
        if MODE_ATTACK:
            return random.randint(1, MAX_COUNTER)

        if len(counters) >= MAX_COUNTER - 1:
            return None

        while True:
            counter = random.randint(1, MAX_COUNTER)
            
            if counter not in counters:
                return counter
        
    def gen_attestation_log(origin):
        now = datetime.now()
        now = now.replace(hour=now.hour, minute=now.minute, second=0, microsecond=0)

        attestation_log = AttestationLog()
        attestation_log.origin = origin
        attestation_log.created_at = now

        attestation_log.counter = AttestationLog.gen_new_counter(origin, now)

        if attestation_log.counter is None:
            return None
        
        return attestation_log
        
    def save_attestation_log(attestation_log):
        session.add(attestation_log)
        session.commit()
    
    def __str__(self):
        return f'{self.origin}@{self.created_at}@{self.counter}' 


Base.metadata.create_all(ENGINE)

if __name__ == '__main__':
    MAX_COUNTER = 5

    attestation_log = AttestationLog.gen_attestation_log('example.com')
    print("--- attestation_log ---\n", attestation_log)

    if attestation_log is not None:
        AttestationLog.save_attestation_log(attestation_log)



