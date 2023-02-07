FROM gomicro/goose

ADD *.sql /migrations/
ADD gooseup.sh /migrations/
RUN chmod -R 777 /migrations/
RUN /migrations/gooseup.sh
