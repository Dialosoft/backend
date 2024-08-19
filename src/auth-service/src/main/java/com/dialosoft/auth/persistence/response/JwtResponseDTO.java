package com.dialosoft.auth.persistence.response;

import lombok.*;

@Getter
@Setter
@AllArgsConstructor
@NoArgsConstructor
@Builder
public class JwtResponseDTO {

    private String accessToken;
    private Long accessTokenExpiresInSeconds;
    private String refreshToken;
    private Long refreshTokenExpiresInSeconds;
}
