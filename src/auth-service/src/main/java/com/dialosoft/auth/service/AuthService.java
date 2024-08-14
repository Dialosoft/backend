package com.dialosoft.auth.service;

import com.dialosoft.auth.persistence.entity.RefreshToken;
import com.dialosoft.auth.persistence.entity.RoleEntity;
import com.dialosoft.auth.persistence.entity.UserEntity;
import com.dialosoft.auth.persistence.repository.RoleRepository;
import com.dialosoft.auth.persistence.repository.UserRepository;
import com.dialosoft.auth.persistence.response.JwtResponseDTO;
import com.dialosoft.auth.persistence.response.ResponseBody;
import com.dialosoft.auth.service.dto.LoginDto;
import com.dialosoft.auth.service.dto.RefreshTokenDto;
import com.dialosoft.auth.service.dto.RegisterDto;
import com.dialosoft.auth.service.utils.RoleType;
import com.dialosoft.auth.web.config.SecurityConfig;
import com.dialosoft.auth.web.config.error.exception.CustomTemplateException;
import com.dialosoft.auth.web.config.jwt.JwtUtil;
import lombok.AllArgsConstructor;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.security.authentication.AuthenticationManager;
import org.springframework.security.authentication.BadCredentialsException;
import org.springframework.security.authentication.UsernamePasswordAuthenticationToken;
import org.springframework.security.core.Authentication;
import org.springframework.security.core.userdetails.UserDetails;
import org.springframework.security.core.userdetails.UsernameNotFoundException;
import org.springframework.stereotype.Service;

import java.util.Collections;
import java.util.Optional;

@Service
@AllArgsConstructor
public class AuthService {

    private final UserRepository userRepository;
    private final RoleRepository roleRepository;
    private final UserSecurityService userSecurityService;
    private final RefreshTokenService refreshTokenService;
    private final TokenBlacklistService tokenBlacklistService;
    private final SecurityConfig securityConfig;
    private final AuthenticationManager authenticationManager;
    private final JwtUtil jwtUtil;

    public ResponseEntity<ResponseBody<?>> register(RegisterDto registerDto) {

        Optional<UserEntity> userEntityOp = userRepository.findByUsernameOrEmail(registerDto.getUsername(), registerDto.getEmail());

        if (userEntityOp.isPresent()) {

            throw new CustomTemplateException("any of the parameters already exists", "Username or email already exists", null, HttpStatus.CONFLICT);
        }

        UserEntity newUser;

        try {
            RoleEntity defaultRole = roleRepository.findByRoleType(RoleType.USER)
                    .orElseThrow(() -> new RuntimeException("Default role not found"));

            UserEntity userEntity = UserEntity.builder()
                    .username(registerDto.getUsername())
                    .email(registerDto.getEmail())
                    .password(securityConfig.encoder().encode(registerDto.getPassword()))
                    .roles(Collections.singleton(defaultRole))
                    .build();

            newUser = userRepository.save(userEntity);
        } catch (Exception e) {

            throw new CustomTemplateException("An error occurred while creating the user", "Internal Server Error", e, HttpStatus.INTERNAL_SERVER_ERROR);
        }

        ResponseBody<?> response = ResponseBody.builder()
                .statusCode(HttpStatus.CREATED.value())
                .message(String.format("User %s created sucessfully", newUser.getId()))
                .metadata(null)
                .build();

        return new ResponseEntity<>(response, HttpStatus.CREATED);
    }

    public ResponseEntity<ResponseBody<?>> login(LoginDto loginDto) {

        UsernamePasswordAuthenticationToken login = new UsernamePasswordAuthenticationToken(
                loginDto.getUsername(),
                loginDto.getPassword()
        );

        try {
            Authentication authentication = this.authenticationManager.authenticate(login);

            if (authentication.isAuthenticated()) {

                String accessToken = jwtUtil.generateAccessToken(loginDto.getUsername(), authentication.getAuthorities());
                Long accessTokenExpiresInSeconds = jwtUtil.getExpirationInSeconds(accessToken);
                RefreshToken refreshToken = refreshTokenService.getOrCreateRefreshTokenByUserName(loginDto.getUsername());

                JwtResponseDTO jwtResponseDTO = JwtResponseDTO.builder()
                        .accessToken(accessToken)
                        .accessTokenExpiresInSeconds(accessTokenExpiresInSeconds)
                        .refreshToken(refreshToken.getRefreshToken())
                        .build();

                ResponseBody<JwtResponseDTO> response = ResponseBody.<JwtResponseDTO>builder()
                        .statusCode(HttpStatus.OK.value())
                        .message("Authentication successfully")
                        .metadata(jwtResponseDTO)
                        .build();

                return new ResponseEntity<>(response, HttpStatus.OK);

            } else {
                throw new UsernameNotFoundException("invalid user request..!!");
            }

        } catch (BadCredentialsException e) {

            throw new CustomTemplateException("Invalid credentials", "Unauthorized", e, HttpStatus.UNAUTHORIZED);
        }
    }

    public ResponseEntity<ResponseBody<JwtResponseDTO>> refreshTokens(RefreshTokenDto refreshTokenDto) {

        return refreshTokenService.findByToken(refreshTokenDto.getRefreshToken())
                .map(refreshTokenService::verifyRefreshTokenExpiration)
                .map(RefreshToken::getUser)
                .map(userInfo -> {

                    String username = userInfo.getUsername();
                    UserDetails userDetails = userSecurityService.loadUserByUsername(username);
                    String accessToken = jwtUtil.generateAccessToken(username, userDetails.getAuthorities());
                    Long accessTokenExpiresInSeconds = jwtUtil.getExpirationInSeconds(accessToken);

                    JwtResponseDTO jwtResponseDTO = JwtResponseDTO.builder()
                            .accessToken(accessToken)
                            .accessTokenExpiresInSeconds(accessTokenExpiresInSeconds)
                            .refreshToken(refreshTokenDto.getRefreshToken()).build();

                    ResponseBody<JwtResponseDTO> response = ResponseBody.<JwtResponseDTO>builder()
                            .statusCode(HttpStatus.OK.value())
                            .message("Pair tokens created successfully")
                            .metadata(jwtResponseDTO)
                            .build();

                    return new ResponseEntity<>(response, HttpStatus.OK);
                })
                .orElseThrow(() -> new RuntimeException(String.format("Refresh token: '%s' was not found in our system", refreshTokenDto.getRefreshToken())));
    }

    public ResponseEntity<ResponseBody<?>> logout(String accessToken) {
        // Get the token user
        String username = jwtUtil.getUsername(accessToken);

        // Find the refresh token associated with the user
        RefreshToken refreshToken = refreshTokenService.getOrCreateRefreshTokenByUserName(username);

        // Save both tokens in the Redis blacklist
        tokenBlacklistService.addToBlacklist(accessToken, jwtUtil.getExpirationInSeconds(accessToken));

        // Delete the refresh token from the database
        refreshTokenService.deleteRefreshTokenByToken(refreshToken.getRefreshToken());

        ResponseBody<?> response = ResponseBody.builder()
                .statusCode(HttpStatus.OK.value())
                .message("Logout successfully")
                .metadata(null)
                .build();

        return new ResponseEntity<>(response, HttpStatus.OK);
    }

}
