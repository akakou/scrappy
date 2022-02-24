from sqlalchemy.ext.declarative import declarative_base
from sqlalchemy import create_engine, Column, Integer, Text, Float, DateTime
from sqlalchemy.orm import scoped_session, sessionmaker
from datetime import datetime
import random, sys, logging

sys.path.append('/attestation/common/')
from wrapper import TERM_LEN

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


class AttestationLog(Base):
    __tablename__ = 'attestation_logs'

    _id = Column('id', Integer, primary_key = True)
    origin = Column('origin', Text)
    created_at = Column(DateTime, default=datetime.now)

    def find_by_origin_and_datetime(origin, datetime):
        attestation_logs = session.query(AttestationLog) \
                .filter(AttestationLog.origin == origin, AttestationLog.created_at == datetime)

        return attestation_logs

    def count_by_origin_and_datetime(origin, datetime):
        attestation_logs = AttestationLog.find_by_origin_and_datetime(origin, datetime)

        return attestation_logs.count()

    def gen_attestation_log(origin):         
        now = datetime.now()
        now = now.replace(hour=now.hour, minute=now.minute, second=now.second // 20, microsecond=0)

        attestation_log = AttestationLog()
        attestation_log.origin = origin
        attestation_log.created_at = now

        if AttestationLog.count_by_origin_and_datetime(origin, now) > 0:
            return None

        return attestation_log
        
    def save_attestation_log(attestation_log):
        session.add(attestation_log)
        session.commit()
    
    def __str__(self):
        timestamp = self.created_at.timestamp()
        return f'{self.origin}@{int(timestamp)}' 


Base.metadata.create_all(ENGINE)

if __name__ == '__main__':

    attestation_log = AttestationLog.gen_attestation_log('example.com')
    print("--- attestation_log ---\n", attestation_log)

    if attestation_log is not None:
        AttestationLog.save_attestation_log(attestation_log)



