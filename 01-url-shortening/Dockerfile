FROM alpine:3.15.11

EXPOSE 8053
ADD go-url-shortner go-url-shortner
RUN chmod +x /go-url-shortner

RUN addgroup -S -g 2000 appgroup && adduser -S appuser -u 1000 -G appgroup
RUN chown -R appuser:appgroup ./go-url-shortner
USER appuser

ENTRYPOINT ["./go-url-shortner"]