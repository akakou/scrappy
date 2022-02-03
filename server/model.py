from sqlalchemy.ext.declarative import declarative_base
from sqlalchemy import create_engine, Column, Integer, Text, Float, DateTime
from sqlalchemy.orm import scoped_session, sessionmaker
from datetime import datetime
import random, sys

ENGINE = create_engine('sqlite:////attestation/ignored_workspace/server_db.sqlite3', echo=True)
Base = declarative_base()

session = scoped_session(
    sessionmaker(
        bind = ENGINE,
        autocommit=False,
    )
)


class AttestationLogForVerifier(Base):
    __tablename__ = 'attestation_logs'

    _id = Column('id', Integer, primary_key = True)
    basename = Column('basename', Text)
    k = Column('k', Text)

    def count(attestation_log):
        attestation_logs = session.query(AttestationLogForVerifier) \
                .filter(AttestationLogForVerifier.basename == attestation_log.basename, AttestationLogForVerifier.k == attestation_log.k).count()

        return attestation_logs

    def exist(attestation_log):
        attestation_logs = AttestationLogForVerifier.count(attestation_log)
        return attestation_logs > 0

    def save(attestation_log):
        session.add(attestation_log)
        session.commit()
    

Base.metadata.create_all(ENGINE)

if __name__ == '__main__':
    attestation_log = AttestationLogForVerifier()
    attestation_log.basename = "basename"
    attestation_log.k = "hello"

    print(AttestationLogForVerifier.exist(attestation_log))

    AttestationLogForVerifier.save(attestation_log)

    print(AttestationLogForVerifier.exist(attestation_log))
